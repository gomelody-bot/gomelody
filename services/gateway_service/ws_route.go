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
	zap.L().Info("New Connection", zap.String("ip", c.RemoteAddr().String()))
	for {
		var msg Message
		err := c.ReadJSON(&msg)
		if err != nil {
			zap.L().Error("failed parsing data", zap.Error(err))
			return
		}

		switch msg.OP {
		case "xyz":
			// call xyz
			break
		}
	}
}
