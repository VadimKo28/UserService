package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewClient(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	const op = "pkg.client.postgres.NewClient"

	db, err := pgxpool.New(ctx, connStr)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: ping failed: %w", op, err)
	}

	return db, nil
}
