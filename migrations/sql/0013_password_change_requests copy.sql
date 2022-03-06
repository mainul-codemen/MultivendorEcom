-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS password_change_requests
(
    id                   Varchar(100) PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id              Varchar(100)     NOT NULL DEFAULT '',
    password_reset_token Varchar(100)     NOT NULL DEFAULT '',
    pass_reset_time      timestamp        DEFAULT current_timestamp
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
