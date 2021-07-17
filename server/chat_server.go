package main

import (
	"context"
	"grpc_chat/chat"
	"io"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const clientIDHeader = "client-id"

type chatServer struct {
	chat.UnimplementedChatServer

	room *room
}

var _ chat.ChatServer = &chatServer{}

func newChatServer() chat.ChatServer {
	return &chatServer{
		room: newRoom(context.Background()),
	}
}

func (s *chatServer) Join(ctx context.Context, req *chat.JoinRequest) (*chat.JoinResponse, error) {
	id := s.room.Join(req.Name)
	return &chat.JoinResponse{Id: id}, nil
}

func (s *chatServer) Leave(ctx context.Context, req *chat.LeaveRequest) (*chat.Empty, error) {
	s.room.Leave(req.Id)
	return &chat.Empty{}, nil
}

func (s *chatServer) Messages(stream chat.Chat_MessagesServer) error {
	fail := make(chan error)

	id, ok := clientID(stream.Context())
	if !ok {
		return status.Errorf(codes.PermissionDenied, "client id required")
	}

	subscription := make(chan chatMessage)
	defer close(subscription)

	s.room.Subscribe(id, subscription)

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return
			} else if err != nil {
				fail <- err
				return
			}

			msg := message{recipientId: in.Id, text: in.Message.Text}
			s.room.Post(msg)
		}
	}()

	go func() {
		for {
			latest := <-subscription

			msg := &chat.ChatMessage{
				Recipient: &chat.Recipient{
					Id:   latest.message.recipientId,
					Name: latest.name,
				},
				Message: &chat.Message{
					Text: latest.message.text,
				},
			}

			err := stream.Send(msg)
			if err != nil {
				fail <- err
			}
		}
	}()

	for {
		select {
		case <-stream.Context().Done():
			s.room.Leave(id)
			return stream.Context().Err()
		case err := <-fail:
			return err
		}
	}
}

func clientID(ctx context.Context) (uint64, bool) {
	m, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(m[clientIDHeader]) == 0 {
		return 0, false
	}
	n, err := strconv.ParseUint(m[clientIDHeader][0], 10, 64)
	if err != nil {
		return 0, false
	}
	return uint64(n), true
}
