package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"user_advt/internal/config"
	"user_advt/internal/domain/users"
	"user_advt/internal/storage"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
	logger *slog.Logger
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

	logger.Info("Storage initialized")

	return &Storage{db: db, logger: logger}, nil
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


func(s *Storage) Get(ctx context.Context, id string) (users.GetUserDTO, error){
	var user users.GetUserDTO

	err := s.db.QueryRow(ctx, "SELECT id, name, email FROM users WHERE id = $1", id).
	  Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
	  if errors.Is(err, sql.ErrNoRows) {
		  s.logger.Error("user not found", slog.String("id", id))
      return users.GetUserDTO{}, storage.ErrUserNotFount 
	  }

		fmt.Printf("execute statement %v:", err)
		return users.GetUserDTO{}, fmt.Errorf("execute statement %w:", err)
	}

  return user, nil
}

func runMigrations(databaseURL string) error {
	const op = "storage.postgres.runMigrations"

	m, err := migrate.New("file:migrations", databaseURL)
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
		"file://./migrations",
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

