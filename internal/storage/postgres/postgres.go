package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"user_advt/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*Storage, error) {
  connStr := cfg.DatabasePath

	// op константа операции для удобства поиска ошибки
	// Используем её в обёртке fmt.Errorf
	const op = "storage.postgres.NewStorage"

	db, err := pgxpool.New(ctx, connStr)

	if err != nil {
		return nil, fmt.Errorf("%s: %w",op, err)
	}

	if err := runMigrations(connStr); err != nil {
		return nil, fmt.Errorf("%s: migration failed: %w", op, err)
	}

	// if err:= rollbackMigrations(connStr); err != nil {
	// 	return nil, fmt.Errorf("%s: rollback migration failed: %w", op, err)
	// }

	return &Storage{db: db}, nil
}

func runMigrations(databaseURL string) error {
	const op = "storage.postgres.runMigrations"

	m, err := migrate.New("file:internal/migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("%s: failed to create migrate instance: %w", op, err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%s: failed to run migrations: %w", op, err)
	}

	return nil
}

func rollbackMigrations(databaseURL string) error {
	m, err := migrate.New(
		"file://./internal/migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	err = m.Down()
	if err == migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}
	if err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	return nil
}

