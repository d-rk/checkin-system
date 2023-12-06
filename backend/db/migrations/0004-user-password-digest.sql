-- +migrate Up
ALTER TABLE users
ADD COLUMN password_digest varchar(255);

ALTER TABLE users
ALTER COLUMN rfid_uid DROP NOT NULL;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

INSERT INTO users (created_at, name) VALUES (current_timestamp, 'admin');
