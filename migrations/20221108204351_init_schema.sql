-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists users
(
    created_at timestamptz default now(),
    updated_at timestamptz,
    deleted_at timestamptz,

    id         uuid        default gen_random_uuid() primary key,
    chat_id    bigint,
    first_name varchar(40),
    last_name  varchar(40),
    wallet_id  uuid        default null,
    language   varchar(2)  default 'en',
    state      int
);

create unique index users_chat_id_idx on users using btree (chat_id);
create unique index users_wallet_idx on users using btree (wallet_id);

create table if not exists providers
(
    created_at timestamptz default now(),
    updated_at timestamptz,
    deleted_at timestamptz,

    id         uuid        default gen_random_uuid() primary key,
    chat_id    bigint,
    user_id    uuid references users (id),

    vcpu       int         default null,
    ram        int         default null,
    storage    int         default null,
    network    int         default null,
    ports      jsonb       default null
);

create unique index providers_chat_id_idx on providers using btree (chat_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table providers;
drop table users;
-- +goose StatementEnd
