package main

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/gomelody-bot/gomelody/pkg/logger"
	"go.uber.org/zap"
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

	// Create new fiber webserver
	app := fiber.New()

	// Register metrics endpoint for prometheus scraping
	prometheus := fiberprometheus.New("voice_service")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// Register websocket endpoint for core interaction
	app.Get("/ws", websocket.New(handleWS))

	// Start fiber server in separate goroutine
	go func() {
		err := app.Listen(cfg.BindAddress)
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

	// Await interruption signal in order to gracefully shutdown webserver
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	zap.L().Info("shutting down...")
}
