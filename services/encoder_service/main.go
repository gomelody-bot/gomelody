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

	// Initialize sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn: cfg.SentryDsn,
	})
	if err != nil {
		zap.L().Error("failed to initialize sentry", zap.Error(err))
	}

	// Initialize new minio client
	min, err := minio.New(cfg.MinioHost, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
	})
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to connect to minio", zap.Error(err))
	}

	// Create new fiber webserver
	app := fiber.New()

	// Register metrics endpoint for prometheus scraping
	prometheus := fiberprometheus.New("voice_service")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// Start fiber server in separate goroutine
	go func() {
		err := app.Listen(cfg.WebAddress)
		if err != nil {
			sentry.CaptureException(err)
			zap.L().Fatal("failed to serve fiber", zap.Error(err))
		}
	}()

	// Defer shutting down of fiber server
	defer func() {
		// Try to gracefully shutdown webserver
		err = app.Shutdown()
		if err != nil {
			zap.L().Fatal("failed to gracefully shutdown webserver", zap.Error(err))
		}
	}()

	// Start TCP server
	l, err := net.Listen("tcp", cfg.BindAddress)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to start tcp server", zap.Error(err))
	}

	// Defer closing of connection after receiving shutdown
	defer func() {
		err := l.Close()
		if err != nil {
			sentry.CaptureException(err)
			zap.L().Error("failed to stop tcp server", zap.Error(err))
		}
	}()

	// Launch new connection handler in separate goroutine
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
	zap.L().Info("shutting down...")
}
