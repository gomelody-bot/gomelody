package server

import (
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"go.uber.org/zap"
)

type Message struct {
	Op   string      `json:"op"`
	Data interface{} `json:"d"`
}

type PlaySoundRequest struct {
	Sound string `json:"sound"`
}

type Server struct {
	app *fiber.App
}

func New() *Server {
	s := &Server{}
	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(s.handleWS))

	s.handleAPI(app.Group("/api"))

	s.app = app
	return s
}

func (s *Server) Start(address string) {
	err := s.app.Listen(address)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed serving websocket", zap.Error(err))
	}
}
