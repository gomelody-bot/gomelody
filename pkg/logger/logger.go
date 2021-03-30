package logger

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"log"
)

func Initialize(dev bool) {
	var l *zap.Logger
	var err error
	if dev {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal("failed to create zap logger: ", err)
		return
	}
	zap.ReplaceGlobals(l)
}
