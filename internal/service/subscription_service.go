package service

import (
	"context"
	"errors"
	"fmt"
	appErrors "subscriptions-api/internal/errors"
	"subscriptions-api/internal/logger"
	"subscriptions-api/internal/model"
	"subscriptions-api/internal/repository"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, sub *model.Subscription) error
	GetSubscription(ctx context.Context, id uint) (*model.Subscription, error)
	ListSubscriptions(ctx context.Context) ([]*model.Subscription, error)
	UpdateSubscription(ctx context.Context, id uint, updateSub *model.UpdateSubscription) error
	DeleteSubscription(ctx context.Context, id uint) error
	CollectStats(ctx context.Context, filter *model.SubscriptionFilter) (*model.SubscriptionStat, error)
}

type subscriptionService struct {
	repository repository.SubscriptionRepository
}

func NewSubscriptionRepository(repository repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repository: repository}
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, sub *model.Subscription) error {
	if sub.EndDate.Before(sub.StartDate) {
		return fmt.Errorf("[SERVICE] End date cannot be earlier than start date")
	}
	err := s.repository.CreateSubscription(ctx, sub)
	if err != nil {
		return fmt.Errorf("[SERVICE] Failed while creating new subscription -> %w", err)
	}
	logger.Info("[SERVICE] OK!")
	return nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id uint) (*model.Subscription, error) {
	sub, err := s.repository.GetSubscription(ctx, id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("[SERVICE] Failed while receiving subscription with id %d -> %w", id, err)
	}
	logger.Info("[SERVICE] OK!")
	return sub, nil
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context) ([]*model.Subscription, error) {
	subs, err := s.repository.ListSubscriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE] Failed while receiving list of subscriptions -> %w", err)
	}
	logger.Info("[SERVICE] OK!")
	return subs, nil
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, id uint, updateSub *model.UpdateSubscription) error {
	if _, err := s.repository.GetSubscription(ctx, id); err != nil {
		return fmt.Errorf("[SERVICE] Failed while receiving subscription -> %w", err)
	}

	err := s.repository.UpdateSubscription(ctx, id, updateSub)
	if err != nil {
		return fmt.Errorf("[SERVICE] Failed while updating subscription with id %d -> %w", id, err)
	}
	logger.Info("[SERVICE] OK!")
	return nil
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, id uint) error {
	err := s.repository.DeleteSubscription(ctx, id)
	if err != nil {
		return fmt.Errorf("[SERVICE] Failed while deleting subscription with id %d -> %w", id, err)
	}
	logger.Info("[SERVICE] OK!")
	return nil
}

func (s *subscriptionService) CollectStats(ctx context.Context, filter *model.SubscriptionFilter) (*model.SubscriptionStat, error) {
	stat, err := s.repository.CollectStats(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE] Failed while collecting stats -> %w", err)
	}
	logger.Info("[SERVICE] OK!")
	return stat, nil
}
