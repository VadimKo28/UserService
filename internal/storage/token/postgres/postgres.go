package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"
	"user_advt/internal/storage"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewTokenStorage(db *pgxpool.Pool, logger *slog.Logger) *Storage {
	return &Storage{db: db, logger: logger}
}

func (s *Storage) Save(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	queryString := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES($1, $2, $3)`

	if _, err := s.db.Exec(ctx, queryString, userID, token, expiresAt); err != nil {
		errorMsg := "insert refresh token failed"

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			s.logger.Error(errorMsg, slog.String("error", pgErr.Detail))
			return storage.ErrInternalServerError
		}

		s.logger.Error(errorMsg, slog.String("error", err.Error()))
		return storage.ErrInternalServerError
	}

	s.logger.Info(queryString)

	return nil
}

func (s *Storage) Get(ctx context.Context, token string) (int, time.Time, error) {
	var userID int
	var expiresAt time.Time

	queryString := `SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1`

	if err := s.db.QueryRow(ctx, queryString, token).Scan(&userID, &expiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Error("refresh token not found", slog.String("error", err.Error()))
			return 0, time.Time{}, storage.ErrRefreshTokenNotFound
		}

		s.logger.Error("fetch refresh token failed", slog.String("error", err.Error()))
		return 0, time.Time{}, storage.ErrInternalServerError
	}

	s.logger.Info(queryString)

	return userID, expiresAt, nil
}

func (s *Storage) Delete(ctx context.Context, token string) error {
	queryString := `DELETE FROM refresh_tokens WHERE token = $1`

	if _, err := s.db.Exec(ctx, queryString, token); err != nil {
		s.logger.Error("delete refresh token failed", slog.String("error", err.Error()))
		return storage.ErrInternalServerError
	}

	s.logger.Info(queryString)

	return nil
}
