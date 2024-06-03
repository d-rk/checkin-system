-- +migrate Up
create table users
(
    id         integer      not null constraint users_pkey primary key,
    created_at timestamp not null,
    updated_at timestamp,
    name       varchar(255) not null constraint users_name_key unique,
    rfid_uid   varchar(255) constraint users_name_rfid unique,
    member_id  varchar(255) constraint users_member_id unique,
    password_digest varchar(255),
    group_name varchar(50),
    role       varchar(50) not null
);

INSERT INTO users (created_at, name, role) VALUES (current_timestamp, 'admin', 'ADMIN');

CREATE INDEX idx_user_groups ON users(group_name);
