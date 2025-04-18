create table if not exists webhooks
(
    id          uuid                     not null
        constraint uni_webhooks_id primary key,
    type        text                     not null,
    label       text                     not null,
    token       text                     not null,
    ignore_host boolean default false    not null,
    created_at  timestamp with time zone not null,
    updated_at  timestamp with time zone not null
);


commit;
