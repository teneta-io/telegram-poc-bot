-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists tasks
(
    created_at   timestamptz                    default now(),
    updated_at   timestamptz,
    deleted_at   timestamptz,

    id           uuid                           default gen_random_uuid() primary key,
    user_id      uuid references users (id),
    vcpu         int                            default null,
    ram          int                            default null,
    storage      int                            default null,
    network      int                            default null,
    image        varchar                        default null,
    ports        jsonb                          default null,
    pub_key      varchar                        default null,

    price_count  float                          default 0,
    price_period varchar                        default 'day',
    status       varchar                        default 'created',

    provider_id  uuid references providers (id) default null,
    remote_user  varchar                        default null,
    remote_ip    varchar                        default null,
    remote_port  varchar                        default null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists tasks;
-- +goose StatementEnd
