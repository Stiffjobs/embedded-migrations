package app

import (
	"context"
	"embed-migrations/internal/database"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	logger     *slog.Logger
	db         *pgxpool.Pool
	router     *http.ServeMux
	migrations fs.FS
}

func New(logger *slog.Logger, migrations fs.FS) *App {

	return &App{
		logger:     logger,
		router:     http.NewServeMux(),
		migrations: migrations,
	}
}

func (a *App) Start(ctx context.Context) error {
	db, err := database.ConnectDB(ctx, a.logger, a.migrations)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	a.db = db

	done := make(chan struct{})

	a.loadRoutes()

	server := http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("failed to listend and serve", slog.Any("error", err))
		}
		close(done)
	}()

	a.logger.Info("Server listening", slog.String("addr", ":8080"))
	select {
	case <-done:
		break
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		server.Shutdown(ctx)
		cancel()
	}

	return nil
}
