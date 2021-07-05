package main

import "grpc_chat/chat"

type chatServer struct {
	chat.UnimplementedChatServer
}

func newServer() chat.ChatServer {
	return &chatServer{}
}
