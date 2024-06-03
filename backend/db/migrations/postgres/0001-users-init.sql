-- +migrate Up
create table users
(
    id         serial       not null constraint users_pkey primary key,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    name       varchar(255) not null constraint users_name_key unique,
    rfid_uid   varchar(255) not null constraint users_name_rfid unique
);
