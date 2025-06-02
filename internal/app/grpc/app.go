package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"chat_service/internal/chat"
	authgrpc "chat_service/internal/grpc/auth"
	chatgrpc "chat_service/internal/grpc/chat"
	"chat_service/internal/kafka"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	log *slog.Logger

	gRPCServer *grpc.Server
	port       int
}

func (a *App) MustRun() {

	if err := a.Run(); err != nil {
		panic(err)
	}

}

func New(log *slog.Logger, authService authgrpc.Auth, chat chat.Chat, counsumer *kafka.KafkaConsumer, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	chatgrpc.Register(gRPCServer, chat, counsumer)

	reflection.Register(gRPCServer)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "grpcapp Run"

	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gprc server is running")

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
