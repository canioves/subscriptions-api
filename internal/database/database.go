package database

import (
	"context"
	"fmt"
	"subscriptions-api/internal/config"
	"subscriptions-api/internal/logger"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context, config *config.Config) (*pgx.Conn, error) {
	logger.Info("[DB] Connecting to database...")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)
	connection, err := pgx.Connect(ctx, dsn)
	if err != nil {
		logger.Error("[DB] Can't connect to database")
		return nil, fmt.Errorf("[DB] Can't connect to database -> %w", err)
	}
	if err := connection.Ping(ctx); err != nil {
		logger.Error("[DB] Ping failed")
		return nil, fmt.Errorf("[DB] Ping failed -> %w", err)
	} else {
		logger.Info("[DB] OK!")
	}
	return connection, nil
}
