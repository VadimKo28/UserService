package postgres

import (
	"app/internal/domain/users"
	"app/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewUserStorage(db *pgxpool.Pool, logger *slog.Logger) *Storage {
	return &Storage{db: db, logger: logger}
}

func (s *Storage) Save(ctx context.Context, user *users.UserCreateDTO) (int, error) {
	var userID int

	queryString := `INSERT INTO users (name, email, password_hash) VALUES($1, $2, $3) RETURNING id`

	if err := s.db.QueryRow(ctx, queryString, user.Name, user.Email, user.Password).Scan(&userID); err != nil {
		errorMsg := "insertion record failure"

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			s.logger.Error(errorMsg, slog.String("error", pgErr.Detail))
			return 0, storage.ErrUserAlreadyExists
		}

		s.logger.Error(errorMsg, slog.String("error", err.Error()))
		return 0, storage.ErrInternalServerError
	}

	s.logger.Info(queryString)

	return userID, nil
}

func (s *Storage) Get(ctx context.Context, id string) (users.User, error) {
	var user users.User

	queryString := `SELECT id, name, email FROM users WHERE id = $1`

	err := s.db.QueryRow(ctx, queryString, id).
		Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Error("user not found", slog.String("id", id))
			return users.User{}, storage.ErrUserNotFount
		}

		fmt.Printf("execute statement %v:", err)
		return users.User{}, fmt.Errorf("execute statement %w:", err)
	}

	s.logger.Info(queryString)

	return user, nil
}

func (s *Storage) GetByCredentials(ctx context.Context, userDTO *users.UserSignInDTO) (users.User, error) {
	var user users.User

	queryString := `SELECT id, name, email FROM users WHERE email = $1 AND password_hash = $2`

	err := s.db.QueryRow(ctx, queryString, userDTO.Email, userDTO.Password).
		Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Error("user not found", slog.String("email", userDTO.Email))
			return users.User{}, storage.ErrUserInvalidCredentials
		}

		fmt.Printf("execute statement %v:", err)
		return users.User{}, fmt.Errorf("execute statement %w:", err)
	}

	s.logger.Info(queryString)

	return user, nil
}
