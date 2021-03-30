package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Dev            bool   `default:"false"`
	WebAddress    string `default:":8020" envconfig:"WEB_ADDRESS"`
	SentryDsn      string `envconfig:"SENTRY_DSN"`
	EncoderService string `envconfig:"ENCODER_SERVICE" required:"true"`
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	var cfg Config
	err := envconfig.Process("API", &cfg)
	if err != nil {
		log.Fatal("failed to load api config: ", err)
	}
	return &cfg
}
