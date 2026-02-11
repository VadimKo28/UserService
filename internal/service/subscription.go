package service

import (
	"app/internal/domain/subscription"
	"context"
	"log/slog"
)

type SubscriptionStorage interface {
	Save(ctx context.Context, subscriptionDTO subscription.CreateSubscriptionDTO) (int, error)
	ListByUserID(ctx context.Context, userID, limit, offset int) ([]subscription.Subscription, error)
	Update(ctx context.Context, subscriptionDTO subscription.UpdateSubscriptionDTO) error
}

type SubscriptionEventPublisher interface {
	Publish(ctx context.Context, event subscription.Event) error
	Close() error
}

type SubscriptionService struct {
	storage   SubscriptionStorage
	publisher SubscriptionEventPublisher
	logger    *slog.Logger
}

func NewSubscriptionService(storage SubscriptionStorage, publisher SubscriptionEventPublisher, logger *slog.Logger) *SubscriptionService {
	return &SubscriptionService{
		storage:   storage,
		publisher: publisher,
		logger:    logger,
	}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, subscriptionDTO subscription.CreateSubscriptionDTO) (int, error) {
	id, err := s.storage.Save(ctx, subscriptionDTO)
	if err != nil {
		return 0, err
	}

	s.publishEvent(ctx, subscription.NewCreatedEvent(id, subscriptionDTO))

	return id, nil
}

func (s *SubscriptionService) GetSubscriptionsByUserID(ctx context.Context, userID, limit, offset int) ([]subscription.Subscription, error) {
	return s.storage.ListByUserID(ctx, userID, limit, offset)
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, subscriptionDTO subscription.UpdateSubscriptionDTO) error {
	if err := s.storage.Update(ctx, subscriptionDTO); err != nil {
		return err
	}

	s.publishEvent(ctx, subscription.NewUpdatedEvent(subscriptionDTO))

	return nil
}

func (s *SubscriptionService) publishEvent(ctx context.Context, event subscription.Event) {
	if s.publisher == nil {
		return
	}

	if err := s.publisher.Publish(ctx, event); err != nil && s.logger != nil {
		s.logger.Error("failed to publish subscription event", slog.String("error", err.Error()))
	}
}
