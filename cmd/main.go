package main

import (
	"chat_service/internal/app"
	"chat_service/internal/config"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Migrations() {
	dbHost := os.Getenv("STORAGE_PATH")
	migrationsPath := "file://./migrations"

	m, err := migrate.New(migrationsPath, dbHost)
	if err != nil {
		panic(fmt.Errorf("failed to create migrate instance: %w", err))
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(fmt.Errorf("failed to apply migrations: %w", err))
	}

}

func main() {
	Migrations()
	slog.Info("Migrations applied successfully")

	cfg := config.MustLoad()

	slog := setupLogger(cfg.Env)

	slog.Info("starting server")

	application := app.New(
		slog, cfg.GPRC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCSrv.Stop()
	slog.Info("server stopped")

}

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
