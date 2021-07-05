package main

import (
	"flag"
	"fmt"
	"grpc_chat/chat"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 20000, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	chat.RegisterChatServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}