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

		// Stream filename to encoder service beforehand
		// TODO: WIP
		_, err = encoder.Write([]byte("super crazy filename\n"))
		if err != nil {
			sentry.CaptureException(err)
			zap.L().Error("failed to provide filename", zap.Error(err))
			return SendError(c, fasthttp.StatusInternalServerError, ErrInternal)
		}

		// Get file from multipart request
		fh, err := c.FormFile("file")
		if err == fasthttp.ErrMissingFile {
			return SendError(c, fasthttp.StatusBadRequest, "upload/missing-file")
		} else if err != nil {
			sentry.CaptureException(err)
			zap.L().Error("failed to access file from multipart body", zap.Error(err))
			return SendError(c, fasthttp.StatusInternalServerError, ErrInternal)
		}

		// Stream file to encoder service
		err = streamFile(fh, encoder)
		if err != nil {
			_ = SendError(c, fasthttp.StatusInternalServerError, ErrInternal)
		}

		// TODO: possibly move up in scope
		// Defer closing of connection until handler is done
		defer func() {
			// Close TCP connection from encoder service
			err := encoder.Close()
			if err != nil {
				sentry.CaptureException(err)
				zap.L().Error("failed to close connection from encoder", zap.Error(err))
				return
			}
			// Set response status code to 200
			err = c.SendStatus(200)
			if err != nil {
				sentry.CaptureException(err)
				zap.L().Error("failed to set response status code to 200", zap.Error(err))
				return
			}
		}()
		return err
	})
}

func streamFile(fh *multipart.FileHeader, c net.Conn) error {
	// Open file from multipart body
	f, err := fh.Open()
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to open file from multipart body", zap.Error(err))
		return err
	}

	// Stream file to encoder service
	_, err = io.Copy(c, f)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to send file to encoder service", zap.Error(err))
		return err
	}
	return nil
}
