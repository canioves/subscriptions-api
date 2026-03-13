package service

import (
	"context"
	"fmt"
	"subscriptions-api/internal/model"
	"subscriptions-api/internal/repository"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, sub *model.Subscription) error
	GetSubscription(ctx context.Context, id uint) (*model.Subscription, error)
}

type subscriptionService struct {
	repository repository.SubscriptionRepository
}

func NewSubscriptionRepository(repository repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repository: repository}
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, sub *model.Subscription) error {
	err := s.repository.CreateSubscription(ctx, sub)
	if err != nil {
		return fmt.Errorf("failed while creating new subscription: %w", err)
	}
	return nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id uint) (*model.Subscription, error) {
	sub, err := s.repository.GetSubscription(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed while getting subscription with id %d: %w", id, err)
	}
	return sub, nil
}
