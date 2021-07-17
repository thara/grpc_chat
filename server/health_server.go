package main

import (
	"context"
	"grpc_chat/chat"
)

type healthServer struct {
	chat.UnimplementedHealthServer
}

func (healthServer) Check(context.Context, *chat.HealthCheckRequest) (*chat.HealthCheckResponse, error) {
	return &chat.HealthCheckResponse{
		Status: chat.HealthCheckResponse_SERVING,
	}, nil
}
