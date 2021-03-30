package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Dev            bool   `default:"false"`
	WebAddress     string `default:":8003" envconfig:"WEB_ADDRESS"`
	BindAddress    string `default:":8004" envconfig:"BIND_ADDRESS"`
	SentryDsn      string `envconfig:"SENTRY_DSN"`
	MinioHost      string `envconfig:"MINIO_HOST" required:"true"`
	MinioAccessKey string `envconfig:"MINIO_ACCESS_KEY" required:"true"`
	MinioSecretKey string `envconfig:"MINIO_SECRET_KEY" required:"true"`
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	var cfg Config
	err := envconfig.Process("ENCODER", &cfg)
	if err != nil {
		log.Fatal("failed to encoder config: ", err)
	}
	return &cfg
}
