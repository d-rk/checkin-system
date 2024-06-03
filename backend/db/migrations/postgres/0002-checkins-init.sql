-- +migrate Up
create table checkins
(
    id                 bigserial not null constraint checkin_pkey primary key,
    date               date not null,
    timestamp          timestamp with time zone,
    user_id            bigint    not null constraint fk_checkins_user references users,
    UNIQUE             (date, user_id)
);

CREATE INDEX idx_checkin_date ON checkins(date);