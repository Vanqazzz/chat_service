package app

import (
	grpcapp "chat_service/internal/app/grpc"
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

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
