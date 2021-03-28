package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
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
			return SendError(c, fasthttp.StatusInternalServerError, ErrInternal)
		}
		_, err = encoder.Write([]byte("super crazy filename\n"))
		if err != nil {
			sentry.CaptureException(err)
			zap.L().Error("could not send filename", zap.Error(err))
			return SendError(c, fasthttp.StatusInternalServerError, ErrInternal)
		}

		// Get file from multipart request
		fh, err := c.FormFile("file")
		if err == fasthttp.ErrMissingFile {
			return SendError(c, fasthttp.StatusBadRequest, "upload/missing-file")
		} else if err != nil {
			sentry.CaptureException(err)
			zap.L().Error("failed to access file", zap.Error(err))
			return SendError(c, fasthttp.StatusInternalServerError, ErrInternal)
		}

		// Stream file to encoder service
		err = streamFile(fh, encoder)
		if err != nil {
			_ = SendError(c, fasthttp.StatusInternalServerError, ErrInternal)
		}


		// Close connection if done
		defer func() {
			err := encoder.Close()
			if err != nil {
				sentry.CaptureException(err)
				zap.L().Error("failed to close connection to encoder", zap.Error(err))
				return
			}
			c.SendStatus(200)
		}()
		return err
	})
}

func streamFile(fh *multipart.FileHeader, c net.Conn) error {
	// Open file
	f, err := fh.Open()
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to open file", zap.Error(err))
		return err
	}

	// Stream to encoder service
	_, err = io.Copy(c, f)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to send file to encoder", zap.Error(err))
		return err
	}
	return nil
}
