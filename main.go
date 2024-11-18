package main

import (
	"context"
	"embed"
	"embed-migrations/internal/app"
	"log/slog"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
)

//go:embed migrations
var migrations embed.FS

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file", slog.Any("error", err))
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	a := app.New(logger, migrations)
	if err := a.Start(ctx); err != nil {
		logger.Error("failed to start server", slog.Any("error", err))
	}

}
