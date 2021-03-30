package main

import (
	"github.com/gofiber/websocket/v2"
	"go.uber.org/zap"
)

type Message struct {
	OP   string      `json:"op"`
	Data interface{} `json:"d"`
}

func handleWS(c *websocket.Conn) {
	// TODO: remove debug log
	zap.L().Info("New Connection", zap.String("ip", c.RemoteAddr().String()))
	for {
		var msg Message
		err := c.ReadJSON(&msg)
		if err != nil {
			zap.L().Error("failed to parse websocket message", zap.Error(err))
			return
		}

		switch msg.OP {
		case "xyz":
			// call xyz
			break
		}
	}
}
