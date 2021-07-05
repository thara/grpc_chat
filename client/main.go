package main

import (
	"context"
	"flag"
	"grpc_chat/chat"
	"log"
	"time"

	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:20000", "The server address in the format of host:port")
	name       = flag.String("name", "Tom", "Your name")
)

func main() {
	flag.Parse()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := chat.NewChatClient(conn)

	var chatName string
	if name == nil {
		chatName = "Unknown"
	} else {
		chatName = *name
	}

	func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		resp, err := client.Join(ctx, &chat.JoinRequest{Name: chatName})
		if err != nil {
			log.Fatalf("%v.Join(_) = _, %v: ", client, err)
		}
		log.Println(resp)
	}()

}
