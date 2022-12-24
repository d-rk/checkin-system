-- +migrate Up
ALTER TABLE users
ADD COLUMN member_id varchar(255) constraint users_member_id unique;