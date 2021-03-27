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

	sentry.Init(sentry.ClientOptions{
		Dsn: cfg.SentryDsn,
	})

	b, err := bot.New(cfg.DiscordToken)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed initializing bot", zap.Error(err))
	}
	b.Start()

	s := server.New()
	go s.Start(cfg.BindAddress)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	zap.L().Info("shutdown...")

	b.Stop()
}

func initializeLogger(dev bool) {
	var (
		l   *zap.Logger
		err error
	)
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
