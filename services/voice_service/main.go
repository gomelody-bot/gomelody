package main

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gomelody-bot/gomelody/pkg/logger"
	"github.com/jonas747/dca"
	"go.uber.org/zap"
	"io"
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

	// TODO: remove debug operations
	encodeSession, err := dca.EncodeFile("test.mp3", dca.StdEncodeOptions)
	if err != nil {
		zap.L().Fatal("failed to encode file", zap.Error(err))
	}
	defer encodeSession.Cleanup()
	output, err := os.Create("output.dca")
	if err != nil {
		zap.L().Fatal("failed to create file", zap.Error(err))
	}
	_, err = io.Copy(output, encodeSession)
	if err != nil {
		zap.L().Fatal("failed to write file", zap.Error(err))
	}

	// Create new fiber webserver
	app := fiber.New()

	// Register metrics endpoint for prometheus scraping
	prometheus := fiberprometheus.New("voice_service")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// Start fiber server in separate goroutine
	go func() {
		err := app.Listen(cfg.BindAddress)
		if err != nil {
			sentry.CaptureException(err)
			zap.L().Fatal("failed to serve fiber", zap.Error(err))
		}
	}()

	// Create new bot instance
	b, err := NewBot(cfg.DiscordToken)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to create bot", zap.Error(err))
	}

	// Start bot instance
	err = b.Start()
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to serve bot", zap.Error(err))
	}

	// Defer stopping of bot instance
	defer func() {
		err = b.Stop()
		if err != nil {
			sentry.CaptureException(err)
			zap.L().Fatal("failed to stop bot gracefully", zap.Error(err))
		}
	}()

	// Await interruption signal in order to gracefully shutdown webserver
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	zap.L().Info("shutting down...")
}
