package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(brokers []string, topic string, groupID string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	return &KafkaConsumer{reader: reader}
}

func (kc *KafkaConsumer) ReadMessages(ctx context.Context, out chan<- kafka.Message) {
	go func() {
		for {
			m, err := kc.reader.ReadMessage(ctx)
			if err != nil {
				close(out)
				return
			}
			out <- m
		}
	}()
}
