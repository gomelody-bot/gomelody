package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Dev          bool   `default:"false"`
	WebAddress  string `default:":8010" envconfig:"WEB_ADDRESS"`
	SentryDsn    string `envconfig:"SENTRY_DSN"`
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	var cfg Config
	err := envconfig.Process("GATEWAY", &cfg)
	if err != nil {
		log.Fatal("failed to load gateway config: ", err)
	}
	return &cfg
}
