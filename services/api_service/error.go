package main

import (
	"encoding/json"
	"errors"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type errorMsg struct {
	Error string `json:"error"`
}

var (
	ErrInternal = "internal"
)

// SendError sends an error message and the according status code to the fiber.Ctx
func SendError(c *fiber.Ctx, status int, error string) error {
	// Set response status code
	err := c.SendStatus(status)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to set response status code", zap.Error(err))
		return errors.New("internal")
	}

	// Encode error message into error struct
	b, err := json.Marshal(errorMsg{Error: error})
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to encode error message", zap.Error(err))
		return errors.New("internal")
	}

	// Write error message to response
	_, err = c.Write(b)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to send error message", zap.Error(err))
		return errors.New("internal")
	}
	return nil
}
