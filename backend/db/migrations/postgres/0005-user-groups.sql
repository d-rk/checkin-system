-- +migrate Up
ALTER TABLE users
ADD COLUMN group_name varchar(50);

ALTER TABLE users
ADD COLUMN role varchar(50) not null default 'USER';

ALTER TABLE users
ALTER COLUMN role drop default;

UPDATE users SET role='ADMIN', updated_at=current_timestamp where name = 'admin';

CREATE INDEX idx_user_groups ON users(group_name);
