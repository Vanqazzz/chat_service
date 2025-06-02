package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

type ChatImpl struct {
	producer *KafkaProducer
}

type MessagePayLoad struct {
	ChatID   string `json:"chat_id"`
	SenderID string `json:"sender_id"`
	Content  string `json:"content"`
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	return &KafkaProducer{writer: writer}
}

func (p *KafkaProducer) SendMessage(ctx context.Context, key string, value []byte) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *ChatImpl) Message(chatID, senderID, content string) (string, error) {

	msg := map[string]string{
		"chat_id":   chatID,
		"sender_id": senderID,
		"content":   content,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	err = p.producer.SendMessage(context.Background(), chatID, data)
	if err != nil {
		return "", err
	}

	return "ok", nil
}
