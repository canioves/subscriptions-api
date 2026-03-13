package repository

import (
	"context"
	"fmt"
	"log"
	"subscriptions-api/internal/model"

	"github.com/jackc/pgx/v5"
)

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, sub *model.Subscription) error
	GetSubscription(ctx context.Context, id uint) (*model.Subscription, error)
	ListSubscriptions(ctx context.Context) ([]*model.Subscription, error)
	UpdateSubscription(ctx context.Context, id uint, sub *model.Subscription) error
	DeleteSubscription(ctx context.Context, id uint) error
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
		subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id uint

	row := r.database.QueryRow(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	)

	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("failed to create subscription with id %d: %w", id, err)
	}
	sub.ID = id
	log.Printf("created new subscription with id %d\n", id)
	return nil
}

func (r *subscriptionRepository) GetSubscription(ctx context.Context, id uint) (*model.Subscription, error) {
	query := "SELECT * FROM subscriptions WHERE id = $1"

	var sub *model.Subscription
	rows, err := r.database.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	sub, err = pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[model.Subscription])

	if err != nil {
		return nil, fmt.Errorf("collect row failed: %w", err)
	}

	log.Printf("got new subscription with id %d\n", id)
	return sub, nil
}

func (r *subscriptionRepository) ListSubscriptions(ctx context.Context) ([]*model.Subscription, error) {
	query := "SELECT * FROM subscriptions"

	var subs []*model.Subscription
	rows, err := r.database.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	subs, err = pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.Subscription])
	if err != nil {
		return nil, fmt.Errorf("collect rows failed: %w", err)
	}
	return subs, nil
}

func (r *subscriptionRepository) UpdateSubscription(ctx context.Context, id uint, sub *model.Subscription) error {
	query := `
	UPDATE subscriptions SET
		service_name = COALESCE($1, service_name),
		price = COALESCE($2, price),
		end_date = COALESCE($3, end_date)
	WHERE id = $4
	`
	row := r.database.QueryRow(ctx, query, sub.ServiceName, sub.Price, sub.EndDate)
	if err := row.Scan(); err != nil {
		return fmt.Errorf("failed to update subscription with id %d: %w", id, err)
	}
	log.Printf("updated subscription with id %d\n", id)
	return nil
}

func (r *subscriptionRepository) DeleteSubscription(ctx context.Context, id uint) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	row := r.database.QueryRow(ctx, query, id)

	if err := row.Scan(); err != nil {
		return fmt.Errorf("failed to delete subscription with id %d: %w", id, err)
	}
	log.Printf("deleted subscription with id %d\n", id)
	return nil
}
