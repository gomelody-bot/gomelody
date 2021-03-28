package main

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gomelody-bot/gomelody/pkg/logger"
	"go.uber.org/zap"
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

	// Start WebServer
	app := fiber.New()
	prometheus := fiberprometheus.New("voice_service")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)
	go func() {
		err := app.Listen(cfg.BindAddress)
		if err != nil {
			sentry.CaptureException(err)
			zap.L().Fatal("failed to serve fiber", zap.Error(err))
		}
	}()

	// Start Bot
	b, err := NewBot(cfg.DiscordToken)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to create bot", zap.Error(err))
	}
	err = b.Start()
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to serve bot", zap.Error(err))
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	zap.L().Info("shutdown...")
}
