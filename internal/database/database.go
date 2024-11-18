package database

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(ctx context.Context, logger *slog.Logger, migrations fs.FS) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(ctx, os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}
	logger.Info("Connected to database")
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("migrate new: %s", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, os.Getenv("DB_URL"))
	if err != nil {
		return nil, fmt.Errorf("migrate new: %s", err)
	}
	logger.Info("Migrator created")

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to migrate db: %w", err)
	}
	logger.Info("Migrations applied")
	return conn, nil
}
