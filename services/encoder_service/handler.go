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
	idGen *snowflake.Node
	min   *minio.Client
}

func NewConnectionHandler(min *minio.Client) *ConnectionHandler {
	node, err := snowflake.NewNode(1)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to initialize snowflake node", zap.Error(err))
	}
	h := &ConnectionHandler{
		idGen: node,
		min:   min,
	}
	return h
}

func (h *ConnectionHandler) handleConnection(c net.Conn) {
	r := bufio.NewReader(c)
	b, err := r.ReadBytes('\n')
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to read filename", zap.Error(err))
		return
	}
	fmt.Println(string(b))

	s, err := dca.EncodeMem(r, dca.StdEncodeOptions)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to encode stream", zap.Error(err))
		return
	}

	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = h.min.PutObject(ctx, "audio-files", fmt.Sprintf("%s.dca", h.idGen.Generate().String()), s, -1, minio.PutObjectOptions{ContentType: "dca"})
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to upload file", zap.Error(err))
	}

	defer ccl()
	defer s.Cleanup()
}
