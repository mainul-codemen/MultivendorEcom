-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS income_tax
(
    id                     Varchar(100) PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id             Varchar(100)     NOT NULL DEFAULT '',
    tax_receipt_number     Varchar(100)     NOT NULL DEFAULT '',
    status                 smallint         DEFAULT 0,
    income_tax_date        DATE NOT NULL,
    tax_amount             double precision NOT NULL DEFAULT 0,
    created_at             TIMESTAMP        DEFAULT current_timestamp,
    created_by             VARCHAR(100)     NOT NULL DEFAULT '',
    updated_at             TIMESTAMP        DEFAULT current_timestamp,
    updated_by             VARCHAR(100)     NOT NULL DEFAULT  '',
    deleted_at             TIMESTAMP        DEFAULT NULL,
    deleted_by             VARCHAR(100)     NOT NULL DEFAULT ''
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
