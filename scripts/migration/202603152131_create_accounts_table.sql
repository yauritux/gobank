-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE USER gobank;
CREATE DATABASE gobank;
GRANT ALL PRIVILEGES ON DATABASE gobank TO gobank;
CREATE TABLE IF NOT EXISTS accounts(
    id uuid not null primary key,
    first_name varchar(100) not null,
    last_name varchar(100) not null,
    account_number varchar(25) not null,
    balance decimal(50, 2) default 0.00
);
ALTER TABLE accounts ADD UNIQUE(account_number);

-- +goose Down
DROP TABLE accounts;