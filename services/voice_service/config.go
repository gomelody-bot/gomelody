package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Dev          bool   `default:"false"`
	BindAddress  string `default:":8000" envconfig:"BIND_ADDRESS"`
	SentryDsn    string `envconfig:"SENTRY_DSN"`
	DiscordToken string `required:"true" envconfig:"DISCORD_TOKEN"`
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	var cfg Config
	err := envconfig.Process("VOICE", &cfg)
	if err != nil {
		log.Fatal("failed loading config: ", err)
	}
	return &cfg
}
