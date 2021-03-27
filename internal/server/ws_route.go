package server

import (
	"github.com/gofiber/websocket/v2"
	"go.uber.org/zap"
)

func (s *Server) handleWS(c *websocket.Conn) {
	zap.L().Info("New Connection", zap.String("ip", c.RemoteAddr().String()))
	for {
		var msg Message
		err := c.ReadJSON(&msg)
		if err != nil {
			zap.L().Error("failed parsing data", zap.Error(err))
			return
		}

		switch msg.Op {
		case "play_sound":
			d, ok := msg.Data.(*PlaySoundRequest)
			if !ok {
				zap.L().Error("invalid play_sound")
				break
			}
			zap.L().Info("Play sound", zap.String("sound", d.Sound))
		}
	}
}
