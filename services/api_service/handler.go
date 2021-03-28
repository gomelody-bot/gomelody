package main

import (
	"errors"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"io"
	"net"
)

type APIHandler struct {
	encoderAddress string
}

func NewAPIHandler(encoderAddress string) *APIHandler {
	return &APIHandler{
		encoderAddress: encoderAddress,
	}
}

func (h *APIHandler) handle(r fiber.Router) {
	r.Post("/upload", func(c *fiber.Ctx) error {
		// Connect to encoder service
		encoder, err := net.Dial("tcp", h.encoderAddress)
		if err != nil {
			sentry.CaptureException(err)
			zap.L().Error("failed to connect to encoder service", zap.Error(err))
			return errors.New("internal")
		}

		fh, err := c.FormFile("file")
		if err == fasthttp.ErrMissingFile {
			return c.SendStatus(fasthttp.StatusBadRequest)
		} else if err != nil {
			sentry.CaptureException(err)
			zap.L().Error("failed to access file", zap.Error(err))
			return err
		}
		f, err := fh.Open()
		if err != nil {
			sentry.CaptureException(err)
			zap.L().Error("failed to open file", zap.Error(err))
			return err
		}
		_, err = io.Copy(encoder, f)
		if err != nil {
			sentry.CaptureException(err)
			zap.L().Error("failed to send file to encoder", zap.Error(err))
		}
		defer func() {
			err := encoder.Close()
			if err != nil {
				sentry.CaptureException(err)
				zap.L().Error("failed to close connection to encoder", zap.Error(err))
			}
		}()
		return err
	})
}
