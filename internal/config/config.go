package config

import (
	"fmt"
	"os"
	"subscriptions-api/internal/logger"

	"github.com/joho/godotenv"
)

type Config struct {
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	AppPort    string
}

func LoadConfig() (*Config, error) {
	logger.Info("[CONFIG] Loading config...")
	envFlag := os.Getenv("GO_ENV")
	config := Config{}

	if envFlag == "docker" {
		config.DBHost = os.Getenv("DB_HOST_DOCKER")
	} else {
		err := godotenv.Load()
		if err != nil {
			logger.Error("[CONFIG] Error with config -> %w", err)
			return nil, fmt.Errorf("[CONFIG] An error occured while loading .env file -> %w", err)
		}

		config.DBHost = os.Getenv("DB_HOST_LOCAL")
	}

	config.DBName = os.Getenv("DB_NAME")
	config.DBUser = os.Getenv("DB_USER")
	config.DBPassword = os.Getenv("DB_PASSWORD")
	config.DBPort = os.Getenv("DB_PORT")
	config.AppPort = os.Getenv("APP_PORT")

	logger.Info("[CONFIG] OK")
	return &config, nil
}
