-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS transaction_source
(
    id                      Varchar(100) PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_source_name VARCHAR(100)     NOT NULL DEFAULT '',
    status                  smallint         DEFAULT 0,
    created_at              TIMESTAMP        DEFAULT current_timestamp,
    created_by              VARCHAR(100)     NOT NULL DEFAULT '',
    updated_at              TIMESTAMP        DEFAULT current_timestamp,
    updated_by              VARCHAR(100)     NOT NULL DEFAULT  '',
    deleted_at              TIMESTAMP        DEFAULT NULL,
    deleted_by              VARCHAR(100)     NOT NULL DEFAULT ''
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
