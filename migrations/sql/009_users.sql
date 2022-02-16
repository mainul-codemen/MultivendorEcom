-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users
(
  id                          VARCHAR(100)     PRIMARY KEY DEFAULT uuid_generate_v4(),
  designation_id              VARCHAR(100)     NOT NULL DEFAULT '',
  user_role                   SMALLINT         DEFAULT 1,        --   COMMENT '1=>Employee. 2=>Developer. 3=>SuperAdmin',
  employee_role               SMALLINT         DEFAULT 1,        --   COMMENT '1=>Employee. 2=>Developer. 3=>SuperAdmin',
  verified_by                 VARCHAR(100)     NOT NULL DEFAULT '',
  join_by                     VARCHAR(100)     NOT NULL DEFAULT '',
  country_id                  VARCHAR(100)     NOT NULL DEFAULT '',
  district_id                 VARCHAR(100)     NOT NULL DEFAULT '',
  station_id                  VARCHAR(100)     NOT NULL DEFAULT '',
  status                      SMALLINT         DEFAULT 1,          --   COMMENT '1=>Inactive. 2=>Active. 3=>Disable. 4=>Leave. 5=>Restrict',
  user_name                   VARCHAR(100)     NOT NULL DEFAULT '',
  first_name                  VARCHAR(100)     NOT NULL DEFAULT '',
  last_name                   VARCHAR(100)     NOT NULL DEFAULT '',
  email                       VARCHAR(100)     NOT NULL DEFAULT '',
  email_verified_at           TIMESTAMP        DEFAULT NULL,
  password                    VARCHAR(100)     NOT NULL DEFAULT '',
  phone_1                     VARCHAR(100)     NOT NULL DEFAULT '',
  phone_2                     VARCHAR(100)     NOT NULL DEFAULT '',
  phone_number_verified_at    SMALLINT         DEFAULT NULL,
  phone_number_verified_code  VARCHAR(100)     NOT NULL DEFAULT '',
  date_of_birth               DATE             DEFAULT NULL,
  gender                      SMALLINT         DEFAULT 1,
  fb_id                       VARCHAR(255)     NOT NULL DEFAULT '',
  photo                       VARCHAR(100)     NOT NULL DEFAULT '',
  nid_front_photo             VARCHAR(100)     NOT NULL DEFAULT '',
  nid_back_photo              VARCHAR(100)     NOT NULL DEFAULT '',
  nid_number                  VARCHAR(100)     NOT NULL DEFAULT '',
  cv_pdf                      VARCHAR(100)     NOT NULL DEFAULT '',
  present_address             VARCHAR(255)     NOT NULL DEFAULT '',
  permanent_address           VARCHAR(255)     NOT NULL DEFAULT '',
  reference                   VARCHAR(100)     NOT NULL DEFAULT '',
  remember_token              VARCHAR(100)     NOT NULL DEFAULT '',
  created_at                  TIMESTAMP        DEFAULT current_timestamp,
  created_by                  VARCHAR(100)     NOT NULL DEFAULT '',
  updated_at                  TIMESTAMP        DEFAULT current_timestamp,
  updated_by                  VARCHAR(100)     NOT NULL DEFAULT  '',
  deleted_at                  TIMESTAMP        DEFAULT NULL,
  deleted_by                  VARCHAR(100)     NOT NULL DEFAULT ''
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd