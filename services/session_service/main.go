package main

import (
	"context"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gomelody-bot/gomelody/pkg/logger"
	"github.com/gomelody-bot/gomelody/pkg/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"time"
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

	// Start fiber server in separate goroutine
	go func() {
		err := app.Listen(cfg.WebAddress)
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

	// Instantiate new redis client
	rd := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB: cfg.RedisDatabase,
	})

	// Initialize context for redis cluster ping execution
	ctx, cl := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cl()

	// Execute ping for verifying successful connection
	sc := rd.Ping(ctx)
	if sc.Err() != nil {
		sentry.CaptureException(sc.Err())
		zap.L().Fatal("failed to ping redis cluster", zap.Error(sc.Err()))
	}

	// Initialize TCP listener for grpc server
	l, err := net.Listen("tcp", cfg.GrpcAddress)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to start listener for grpc server", zap.Error(err))
	}

	// Start grpc server and serve on above TCP listener
	grpcServer := grpc.NewServer()
	proto.RegisterSessionServiceServer(grpcServer, NewSessionServer(rd))
	err = grpcServer.Serve(l)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to serve grpc server", zap.Error(err))
	}

	zap.L().Info("shutting down...")
}
