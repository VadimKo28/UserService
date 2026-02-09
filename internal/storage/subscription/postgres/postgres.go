package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"user_advt/internal/domain/subscription"
	"user_advt/internal/storage"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewSubscriptionStorage(db *pgxpool.Pool, logger *slog.Logger) *Storage {
	return &Storage{db: db, logger: logger}
}

func (s *Storage) Save(ctx context.Context, subscriptionDto subscription.CreateSubscriptionDTO) (int, error) {
	var id int
	queryString := `INSERT INTO subscriptions (user_id, service_name, price, start_date, end_date) VALUES($1, $2, $3, $4, $5) RETURNING id`

	endDate := any(subscriptionDto.EndDate)
	if subscriptionDto.EndDate.IsZero() {
		endDate = nil
	}

	if err := s.db.QueryRow(ctx,
		queryString,
		subscriptionDto.UserID,
		subscriptionDto.ServiceName,
		subscriptionDto.Price,
		subscriptionDto.StartDate,
		endDate,
	).Scan(&id); err != nil {

		errorMsg := "insertion record failure"

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				// TODO Сейчас ошибка 23505 никогда не случится при создании подписки,
				// т.к. это unique_violation, такого ограничения для подписок в базе нет,
				// можно убрать эту ветку с проверкой,
				// но пока просто оставлю на будущее
				s.logger.Error(errorMsg, slog.String("error", pgErr.Detail))
				return 0, storage.ErrInternalServerError
			}
			if pgErr.Code == "23503" {
				s.logger.Error(errorMsg, slog.String("error", pgErr.Detail))
				return 0, storage.ErrUserNotFount
			}
		}

		s.logger.Error(errorMsg, slog.String("error", err.Error()))
		return 0, storage.ErrInternalServerError
	}

	s.logger.Info(queryString)

	return id, nil
}

func (s *Storage) ListByUserID(ctx context.Context, userID, limit, offset int) ([]subscription.Subscription, error) {
	queryString := `SELECT id, user_id, service_name, price, start_date, end_date FROM subscriptions WHERE user_id = $1 ORDER BY id LIMIT $2 OFFSET $3`

	rows, err := s.db.Query(ctx, queryString, userID, limit, offset)
	if err != nil {
		s.logger.Error("failed to load subscriptions", slog.String("error", err.Error()))
		return nil, storage.ErrInternalServerError
	}
	defer rows.Close()

	subscriptions := make([]subscription.Subscription, 0)
	for rows.Next() {
		var item subscription.Subscription
		var endDate sql.NullTime

		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.ServiceName,
			&item.Price,
			&item.StartDate,
			&endDate,
		); err != nil {
			s.logger.Error("failed to scan subscription row", slog.String("error", err.Error()))
			return nil, storage.ErrInternalServerError
		}

		if endDate.Valid {
			item.EndDate = endDate.Time
		}

		subscriptions = append(subscriptions, item)
	}

	if err := rows.Err(); err != nil {
		s.logger.Error("failed to iterate subscriptions", slog.String("error", err.Error()))
		return nil, storage.ErrInternalServerError
	}

	return subscriptions, nil
}

func (s *Storage) Update(ctx context.Context, subscriptionDto subscription.UpdateSubscriptionDTO) error {
	queryString := `UPDATE subscriptions SET service_name = $1, price = $2, start_date = $3, end_date = $4 WHERE id = $5 AND user_id = $6`

	endDate := any(subscriptionDto.EndDate)
	if subscriptionDto.EndDate.IsZero() {
		endDate = nil
	}

	result, err := s.db.Exec(
		ctx,
		queryString,
		subscriptionDto.ServiceName,
		subscriptionDto.Price,
		subscriptionDto.StartDate,
		endDate,
		subscriptionDto.ID,
		subscriptionDto.UserID,
	)
	if err != nil {
		s.logger.Error("failed to update subscription", slog.String("error", err.Error()))
		return storage.ErrInternalServerError
	}

	if result.RowsAffected() == 0 {
		return storage.ErrSubscriptionNotFound
	}

	return nil
}
