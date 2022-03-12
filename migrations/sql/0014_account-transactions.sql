-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS accounts_transactions
(
    id                     Varchar(255) PRIMARY KEY DEFAULT uuid_generate_v4(),
    from_account_id        Varchar(100)     NOT NULL DEFAULT '',
    to_account_id          Varchar(100)     NOT NULL DEFAULT '',
    user_id                Varchar(100)     NOT NULL DEFAULT '',
    transaction_amount     double precision NOT NULL DEFAULT 0,
    transaction_type_id    Varchar(100)     NOT NULL DEFAULT '',
    transaction_source_id  Varchar(100)     NOT NULL DEFAULT '',
    reference              Varchar(100)     NOT NULL DEFAULT '',
    note                   Varchar(100)     NOT NULL DEFAULT '',
    status                 smallint         DEFAULT 0,
    from_acnt_previous_balance       double precision NOT NULL DEFAULT 0,
    from_acnt_current_balance        double precision NOT NULL DEFAULT 0,
    to_acnt_current_balance        double precision NOT NULL DEFAULT 0,
    to_acnt_current_balance        double precision NOT NULL DEFAULT 0,
    created_at             TIMESTAMP        DEFAULT current_timestamp,
    created_by             VARCHAR(100)     NOT NULL DEFAULT '',
    accepted_at            TIMESTAMP        DEFAULT current_timestamp,
    accepted_by            VARCHAR(100)     NOT NULL DEFAULT '',
    updated_at             TIMESTAMP        DEFAULT current_timestamp,
    updated_by             VARCHAR(100)     NOT NULL DEFAULT  '',
    deleted_at             TIMESTAMP        DEFAULT NULL,
    deleted_by             VARCHAR(100)     NOT NULL DEFAULT ''
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
