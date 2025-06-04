package suite

import (
	"chat_service/internal/app"
	"chat_service/internal/config"
	"context"
	"log/slog"
	"net"
	"os"
	"strconv"
	"testing"

	protos "github.com/Vanqazzz/protos/gen/go/chat_service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient protos.AuthClient
}

const (
	grpcHost = "localhost"
)

// Starting Test gRPC Server
func StartTestGRPCServer(cfg *config.Config) {

	slog := setupLogger(cfg.Env)

	slog.Info("starting server")

	application := app.New(slog, cfg.GPRC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/local_tests.yaml")
	StartTestGRPCServer(cfg)

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GPRC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(grpcAddress(cfg), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: protos.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GPRC.Port))
}

// Test Logger
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return log
}
