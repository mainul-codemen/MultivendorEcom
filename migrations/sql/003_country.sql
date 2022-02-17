-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS country
(
    id          Varchar(100) PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        varchar(100)     NOT NULL DEFAULT '',
    status      smallint         DEFAULT 0,
    position    int              DEFAULT 0,
    created_at  timestamp        DEFAULT current_timestamp,
    created_by  varchar(100)     NOT NULL DEFAULT '',
    updated_at  timestamp        DEFAULT current_timestamp,
    updated_by  varchar(100)     NOT NULL DEFAULT  '',
    deleted_at  timestamp        DEFAULT NULL,
    deleted_by  varchar(100)     NOT NULL DEFAULT ''
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
