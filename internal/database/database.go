package database

import (
	"context"
	"fmt"
	"log"
	"subscriptions-api/internal/config"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context, config *config.Config) (*pgx.Conn, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)
	connection, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}
	if err := connection.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	} else {
		log.Println("succesfully connected to database")
	}
	return connection, nil
}
