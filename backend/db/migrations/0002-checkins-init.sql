-- +migrate Up
create table checkins
(
    id                 bigserial not null constraint checkin_pkey primary key,
    timestamp          timestamp with time zone,
    user_id            bigint    not null constraint fk_checkins_user references users
);
