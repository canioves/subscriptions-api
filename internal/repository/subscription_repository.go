package repository

import (
	"context"
	"errors"
	"fmt"
	appErrors "subscriptions-api/internal/errors"
	"subscriptions-api/internal/logger"
	"subscriptions-api/internal/model"

	"github.com/jackc/pgx/v5"
)

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, sub *model.Subscription) error
	GetSubscription(ctx context.Context, id uint) (*model.Subscription, error)
	ListSubscriptions(ctx context.Context) ([]*model.Subscription, error)
	UpdateSubscription(ctx context.Context, id uint, updateSub *model.UpdateSubscription) error
	DeleteSubscription(ctx context.Context, id uint) error
	CollectStats(ctx context.Context, filter *model.SubscriptionFilter) (*model.SubscriptionStat, error)
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
		return fmt.Errorf("[REPO] Query failed -> %w", err)
	}
	sub.ID = id
	logger.Info("[REPO] OK!")
	return nil
}

func (r *subscriptionRepository) GetSubscription(ctx context.Context, id uint) (*model.Subscription, error) {
	query := "SELECT * FROM subscriptions WHERE id = $1"

	var sub *model.Subscription
	rows, err := r.database.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("[REPO] Query failed -> %w", err)
	}

	sub, err = pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[model.Subscription])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, fmt.Errorf("[REPO] Collect row failed -> %w", err)
	}

	logger.Info("[REPO] OK!")
	return sub, nil
}

func (r *subscriptionRepository) ListSubscriptions(ctx context.Context) ([]*model.Subscription, error) {
	query := "SELECT * FROM subscriptions"

	var subs []*model.Subscription
	rows, err := r.database.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("[REPO] Query failed -> %w", err)
	}

	subs, err = pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.Subscription])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, fmt.Errorf("[REPO] Collect rows failed -> %w", err)
	}
	logger.Info("[REPO] OK!")
	return subs, nil
}

func (r *subscriptionRepository) UpdateSubscription(ctx context.Context, id uint, updateSub *model.UpdateSubscription) error {
	query := `
	UPDATE subscriptions SET
		service_name = COALESCE($1, service_name),
		price = COALESCE($2, price),
		end_date = COALESCE($3, end_date)
	WHERE id = $4
	`
	_, err := r.database.Exec(ctx, query, updateSub.ServiceName, updateSub.Price, updateSub.EndDate, id)
	if err != nil {
		return fmt.Errorf("[REPO] Failed to execute query -> %w", err)
	}
	updateSub.ID = id

	logger.Info("[REPO] OK!")
	return nil
}

func (r *subscriptionRepository) DeleteSubscription(ctx context.Context, id uint) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := r.database.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("[REPO] Failed to execute query -> %w", err)
	}

	logger.Info("[REPO] OK!")
	return nil
}

func (r *subscriptionRepository) CollectStats(ctx context.Context, filter *model.SubscriptionFilter) (*model.SubscriptionStat, error) {
	query := `
	SELECT
		COUNT(id) as total_count,
		COALESCE(SUM(price), 0) as total_sum
	FROM subscriptions
	WHERE (start_date <= $2 and (end_date >= $1 or end_date is null))
	`

	if filter.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name ILIKE '%s'", *filter.ServiceName)
	}

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = '%s'", *filter.UserID)
	}

	var stat *model.SubscriptionStat

	rows, err := r.database.Query(ctx, query, filter.StartPeriod, filter.EndPeriod)

	if err != nil {
		return nil, fmt.Errorf("[REPO] Query failed -> %w", err)
	}

	stat, err = pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[model.SubscriptionStat])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, fmt.Errorf("[REPO] Collect rows failed -> %w", err)
	}

	logger.Info("[REPO] OK!")
	return stat, nil
}
