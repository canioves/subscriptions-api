-- +goose Up
CREATE TABLE "subscriptions" (
    id SERIAL PRIMARY KEY,
    service_name TEXT NOT NULL,
    price INT NOT NULL,
    user_id UUID NOT NULL,
    start_date timestamptz NOT NULL,
    end_date timestamptz
);

-- +goose Down
DROP TABLE "subscriptions" CASCADE;
