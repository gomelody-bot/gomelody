package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/getsentry/sentry-go"
	"github.com/jonas747/dca"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"net"
	"time"
)

type ConnectionHandler struct {
	gen *snowflake.Node
	min *minio.Client
}

func NewConnectionHandler(min *minio.Client) *ConnectionHandler {
	node, err := snowflake.NewNode(1)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to initialize snowflake node", zap.Error(err))
	}
	h := &ConnectionHandler{
		gen: node,
		min: min,
	}
	return h
}

func (h *ConnectionHandler) handleConnection(c net.Conn) {
	r := bufio.NewReader(c)

	// Read filename from byte array until delimiter
	b, err := r.ReadBytes('\n')
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to read filename", zap.Error(err))
		return
	}

	// TODO: remove debug log
	fmt.Println(string(b))

	// Initialize encoding session for file upload
	s, err := dca.EncodeMem(r, dca.StdEncodeOptions)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to encode stream", zap.Error(err))
		return
	}

	// Initialize file upload to minio
	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = h.min.PutObject(ctx, "audio-files", fmt.Sprintf("%s.dca", h.gen.Generate().String()), s, -1, minio.PutObjectOptions{ContentType: "dca"})
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to upload file to minio", zap.Error(err))
	}

	// Cancel context in order to release its associated resources
	defer ccl()

	// Clean up encoding session
	defer s.Cleanup()
}
