package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gomelody-bot/gomelody/pkg/proto"
	"go.uber.org/zap"
)

// TODO: implement LRU mechanism for sessions (max 10 e.g.)

type server struct {
	rd *redis.Client
	proto.UnimplementedSessionServiceServer
}

func NewSessionServer(rd *redis.Client) *server {
	return &server{
		rd: rd,
	}
}

func (s server) CreateSession(ctx context.Context, request *proto.CreateSessionRequest) (*proto.CreateSessionResponse, error) {
	zap.L().Info("got request", zap.Any("request", request))
	return &proto.CreateSessionResponse{}, nil
}

func (s server) DeleteSession(ctx context.Context, request *proto.DeleteSessionRequest) (*proto.DeleteSessionResponse, error) {
	panic("implement me")
}

func (s server) RefreshSession(ctx context.Context, request *proto.RefreshSessionRequest) (*proto.RefreshSessionResponse, error) {
	panic("implement me")
}

func (s server) VerifySession(ctx context.Context, request *proto.VerifySessionRequest) (*proto.VerifySessionResponse, error) {
	panic("implement me")
}

func (s server) GetSessions(ctx context.Context, request *proto.GetSessionsRequest) (*proto.GetSessionsResponse, error) {
	panic("implement me")
}
