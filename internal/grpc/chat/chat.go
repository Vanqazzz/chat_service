package chat

import (
	"chat_service/internal/chat"
	"chat_service/internal/kafka"
	"context"
	"encoding/json"
	"fmt"

	protos "github.com/Vanqazzz/protos/gen/go/chat_service/chat"
	segmentioKafka "github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

type serverAPI struct {
	protos.UnimplementedChatServiceServer
	chat chat.Chat

	consumer *kafka.KafkaConsumer
}

func Register(gRPC *grpc.Server, chat chat.Chat, consumer *kafka.KafkaConsumer) {
	protos.RegisterChatServiceServer(gRPC, &serverAPI{chat: chat, consumer: consumer})
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

func (s *serverAPI) ReceiveMessages(req *protos.UserStreamRequest, stream protos.ChatService_ReceiveMessagesServer) error {
	ctx := stream.Context()

	msgChan := make(chan segmentioKafka.Message)

	go s.consumer.ReadMessages(ctx, msgChan)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}
			var parsed kafka.MessagePayLoad
			err := json.Unmarshal(msg.Value, &parsed)
			if err != nil {
				continue
			}

			if err := stream.Send(&protos.Message{
				ChatId:   parsed.ChatID,
				SenderId: parsed.SenderID,
				Content:  parsed.Content,
			}); err != nil {
				return err
			}

		}
	}

}
