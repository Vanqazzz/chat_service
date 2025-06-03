package app

import (
	grpcapp "chat_service/internal/app/grpc"
	"chat_service/internal/chat"
	kf "chat_service/internal/kafka"
	"chat_service/internal/services/auth"
	"chat_service/internal/storage"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := storage.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	brokers := []string{"kafka:9092"}
	topic := "chat-messages"
	groupID := "chat-group"

	producer := kf.NewKafkaProducer(brokers, topic)
	consumer := kf.NewKafkaConsumer(brokers, topic, groupID)

	chatService := chat.NewChatService(producer)
	grpcApp := grpcapp.New(log, authService, chatService, consumer, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
