package main

import (
	"context"
	"github.com/kushturner/tfl-alerts/internal/app"
	"log/slog"
	"os"
	"os/signal"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	a := app.New(logger)

	if err := a.Start(ctx); err != nil {
		logger.Error("failed to start server", slog.Any("error", err))
	}
}
