package main

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gomelody-bot/gomelody/pkg/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
	"net"
	"os"
	"os/signal"
)

func main() {
	cfg := LoadConfig()
	logger.Initialize(cfg.Dev)

	// Initialize Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn: cfg.SentryDsn,
	})
	if err != nil {
		zap.L().Error("failed to initialize sentry", zap.Error(err))
	}

	// Initialize MinIO
	min, err := minio.New(cfg.MinioHost, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
	})
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to connect to MinIO", zap.Error(err))
	}

	// Start WebServer
	app := fiber.New()
	prometheus := fiberprometheus.New("voice_service")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)
	go func() {
		err := app.Listen(cfg.WebAddress)
		if err != nil {
			sentry.CaptureException(err)
			zap.L().Fatal("failed to serve fiber", zap.Error(err))
		}
	}()

	// Start TCP Server
	l, err := net.Listen("tcp", cfg.BindAddress)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to start tcp server", zap.Error(err))
	}
	defer func() {
		err := l.Close()
		if err != nil {
			sentry.CaptureException(err)
			zap.L().Error("failed to stop tcp server", zap.Error(err))
		}
	}()

	// Handle TCP Connection
	go func() {
		h := NewConnectionHandler(min)
		for {
			c, err := l.Accept()
			if err != nil {
				sentry.CaptureException(err)
				zap.L().Error("failed to accept connection", zap.Error(err))
				return
			}
			h.handleConnection(c)
		}
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	zap.L().Info("shutdown initialized...")
}
