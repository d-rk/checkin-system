-- +migrate Up
ALTER TABLE users
ADD COLUMN password_hash varchar(255) not null default '-';
