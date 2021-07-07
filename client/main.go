package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"grpc_chat/chat"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

	id, myName := join(client, name)
	fmt.Println("client ID: ", id)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "client-id", strconv.FormatUint(id, 10))
	stream, err := client.Messages(ctx)
	if err != nil {
		log.Fatalf("%v.Messages(_) = _, %v: ", client, err)
	}

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return
			} else if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			fmt.Printf("\033[0G%s > %s\n", in.Recipient.Name, in.Message.Text)
			fmt.Printf("%s > ", myName)
		}
	}()

	fmt.Printf("%s > ", myName)

	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		text := stdin.Text()

		if text == "exit" {
			leave(client, id)
			fmt.Println("bye.")
			return
		}

		req := &chat.MessageRequest{
			Id: id,
			Message: &chat.Message{
				Text: text,
			},
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}

		fmt.Printf("%s > ", myName)
	}
}

func join(client chat.ChatClient, name *string) (uint64, string) {
	var chatName string
	if name == nil {
		chatName = "Unknown"
	} else {
		chatName = *name
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := client.Join(ctx, &chat.JoinRequest{Name: chatName})
	if err != nil {
		log.Fatalf("%v.Join(_) = _, %v: ", client, err)
	}

	return resp.Id, chatName
}

func leave(client chat.ChatClient, id uint64) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := client.Leave(ctx, &chat.LeaveRequest{Id: id})
	if err != nil {
		log.Fatalf("%v.Leave(_) = _, %v: ", client, err)
	}
}
