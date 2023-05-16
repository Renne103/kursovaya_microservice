-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_balance (
     user_id BIGINT NOT NULL UNIQUE,
     balance NUMERIC(6, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS balance_holder (
     order_id BIGINT NOT NULL UNIQUE,
     service_id BIGINT NOT NULL,
     user_id BIGINT NOT NULL,
     service_name TEXT NOT NULL,
     cost NUMERIC(6, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS service_history (
     order_id BIGINT NOT NULL UNIQUE,
     service_id BIGINT NOT NULL,
     user_id BIGINT NOT NULL,
     operation_time DATE,
     service_name TEXT NOT NULL,
     cost NUMERIC(6, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS transaction_history (
     transaction_id BIGINT NOT NULL UNIQUE,
     user_id BIGINT NOT NULL,
     operation_time TIMESTAMP,
     description TEXT NOT NULL,
     cost NUMERIC(6, 2) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS pagination_time
    ON transaction_history (operation_time, transaction_id);

CREATE UNIQUE INDEX IF NOT EXISTS pagination_cost
    ON transaction_history (cost, transaction_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_balance;

DROP TABLE IF EXISTS balance_holder;

DROP TABLE IF EXISTS service_history;

DROP TABLE IF EXISTS transaction_history;
DROP INDEX IF EXISTS pagination_cost;
DROP INDEX IF EXISTS pagination_time;
-- +goose StatementEnd
