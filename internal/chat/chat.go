package chat

import (
	"chat_service/internal/kafka"
	"context"
	"encoding/json"
	"fmt"
)

type Chat interface {
	Message(
		chat_id string,
		senderID string,
		content string,
	) (string, error)
}

type chatService struct {
	producer *kafka.KafkaProducer
}

func NewChatService(producer *kafka.KafkaProducer) Chat {
	return &chatService{producer: producer}
}

func (c *chatService) Message(chatID string, senderID string, content string) (string, error) {

	payload := kafka.MessagePayLoad{
		ChatID:   chatID,
		SenderID: senderID,
		Content:  content,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal message payload: %w", err)

	}

	err = c.producer.SendMessage(context.Background(), chatID, data)
	if err != nil {
		return "", fmt.Errorf("failed to send message to kafka: %w", err)
	}

	return "message sent", nil
}
