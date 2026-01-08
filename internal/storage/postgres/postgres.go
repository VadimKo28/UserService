package postgres

import (
	"context"
	"fmt"
	"user_advt/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(ctx context.Context, cfg *config.Config) (*Storage, error) {
  connStr := cfg.DatabasePath

	db, err := pgxpool.New(ctx, connStr)

	if err != nil {
		return nil, fmt.Errorf("postgres pool error: %w", err)
	}

	stmt := 
		`CREATE TABLE IF NOT EXISTS urls (
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias);`
	
	_, err = db.Exec(ctx, stmt) 

	if err != nil {
		return nil, fmt.Errorf("create table error: %w", err)
	}

	fmt.Println("Database init")

	return &Storage{db: db}, nil
}
