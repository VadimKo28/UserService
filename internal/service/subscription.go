package service

import (
	"context"
	"user_advt/internal/domain/subscription"
)

type SubscriptionStorage interface {
	Save(ctx context.Context, subscriptionDTO subscription.CreateSubscriptionDTO) (int, error)
	ListByUserID(ctx context.Context, userID, limit, offset int) ([]subscription.Subscription, error)
	Update(ctx context.Context, subscriptionDTO subscription.UpdateSubscriptionDTO) error
}

type SubscriptionService struct {
	storage SubscriptionStorage
}

func NewSubscriptionService(storage SubscriptionStorage) *SubscriptionService {
	return &SubscriptionService{storage: storage}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, subscriptionDTO subscription.CreateSubscriptionDTO) (int, error) {
	return s.storage.Save(ctx, subscriptionDTO)
}

func (s *SubscriptionService) GetSubscriptionsByUserID(ctx context.Context, userID, limit, offset int) ([]subscription.Subscription, error) {
	return s.storage.ListByUserID(ctx, userID, limit, offset)
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, subscriptionDTO subscription.UpdateSubscriptionDTO) error {
	return s.storage.Update(ctx, subscriptionDTO)
}
