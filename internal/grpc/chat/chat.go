package chat

import (
	"context"
	"fmt"

	protos "github.com/Vanqazzz/protos/gen/go/chat_service/chat"
	"google.golang.org/grpc"
)

type Chat interface {
	Message(
		chat_id string,
		sender_id string,
		content string,
	) (status string, err error)
}

type serverAPI struct {
	protos.UnimplementedChatServiceServer
	chat Chat
}

func Register(gRPC *grpc.Server, chat Chat) {
	protos.RegisterChatServiceServer(gRPC, &serverAPI{chat: chat})
}

func (s *serverAPI) SendMessage(ctx context.Context, req *protos.MessageRequest) (*protos.MessageResponse, error) {

	status, err := s.chat.Message(req.ChatId, req.SenderId, req.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return &protos.MessageResponse{
		Status: status,
	}, nil

}

/* func (s *serverAPI) ReceiveMessage(req *protos.UserStreamRequest, stream protos.ChatService_ReceiveMessagesServer) error {

} */
