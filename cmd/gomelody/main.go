package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"

	"github.com/gomelody-bot/gomelody/internal/bot"
	"github.com/gomelody-bot/gomelody/internal/config"
	"github.com/gomelody-bot/gomelody/internal/server"
)

func main() {
	cfg := config.LoadConfig()
	initializeLogger(cfg.Dev)

	err := sentry.Init(sentry.ClientOptions{
		Dsn: cfg.SentryDsn,
	})
	if err != nil {
		zap.L().Error("failed to initialize sentry", zap.Error(err))
	}

	b, err := bot.New(cfg.DiscordToken)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed initializing bot", zap.Error(err))
	}
	err = b.Start()
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to connect to to Discord", zap.Error(err))
	}

	s := server.New()
	go s.Start(cfg.BindAddress)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	zap.L().Info("shutdown...")

	err = b.Stop()
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to close connection", zap.Error(err))
	}
}

func initializeLogger(dev bool) {
	var l *zap.Logger
	var err error
	if dev {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal("failed creating logger: ", err)
		return
	}
	zap.ReplaceGlobals(l)
}
