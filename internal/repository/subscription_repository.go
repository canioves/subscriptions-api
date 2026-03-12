package repository

import (
	"context"
	"fmt"
	"log"
	"subscriptions-api/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, sub *model.Subscription) error
	GetSubscription(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	ListSubscriptions(ctx context.Context) (*[]model.Subscription, error)
	UpdateSubscription(ctx context.Context, sub *model.Subscription) error
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
}

type subscriptionRepository struct {
	database *pgx.Conn
}

func NewSubscriptionRepository(database *pgx.Conn) SubscriptionRepository {
	return &subscriptionRepository{database: database}
}

func (r *subscriptionRepository) CreateSubscription(ctx context.Context, sub *model.Subscription) error {
	query := `
		INSERT INTO 
		subscriptions (id, service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var id uint

	row := r.database.QueryRow(
		ctx,
		query,
		sub.ID,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	)

	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("failed to create subscription with id %s: %w", id, err)
	}

	log.Println("created new subscription with id %d", id)

	return nil
}
