-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS delivery_charge
(
    id              varchar(100)     PRIMARY KEY DEFAULT uuid_generate_v4(),
    country_id      varchar(100)     NOT NULL DEFAULT '',
    district_id     varchar(100)     NOT NULL DEFAULT '',
    station_id      varchar(100)     NOT NULL DEFAULT '',
    weight_min      double precision NOT NULL DEFAULT 0,
    weight_max      double precision NOT NULL DEFAULT 0,
    delivery_charge double precision NOT NULL DEFAULT 0,
    dc_status       smallint         DEFAULT 0,
    created_at      timestamp        DEFAULT current_timestamp,
    created_by      varchar(100)     NOT NULL DEFAULT '',
    updated_at      timestamp        DEFAULT current_timestamp,
    updated_by      varchar(100)     NOT NULL DEFAULT  '',
    deleted_at      timestamp        DEFAULT NULL,
    deleted_by      varchar(100)     NOT NULL DEFAULT ''
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
