#!/bin/sh
set -e

echo "DB_HOST: $DB_HOST"
echo "DB_USER: $DB_USER"

DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST_DOCKER}:${DB_PORT}/${DB_NAME}?sslmode=disable"

echo "Running migrations..."
goose -dir ./migrations postgres "$DB_URL" up

exec ./server