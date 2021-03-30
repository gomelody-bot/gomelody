package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Dev           bool   `default:"false"`
	WebAddress    string `default:":8040" envconfig:"WEB_ADDRESS"`
	GrpcAddress   string `default:":8041" envconfig:"GRPC_ADDRESS"`
	SentryDsn     string `envconfig:"SENTRY_DSN"`
	RedisAddr     string `default:"localhost:6379" envconfig:"REDIS_ADDRESS"`
	RedisPassword string `default:"" envconfig:"REDIS_PASSWORD"`
	RedisDatabase int    `default:"0" envconfig:"REDIS_DATABASE"`
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	var cfg Config
	err := envconfig.Process("SESSION", &cfg)
	if err != nil {
		log.Fatal("failed to load session config: ", err)
	}
	return &cfg
}
