create table if not exists updates
(
    id          uuid                     not null
        constraint uni_updates_id primary key,
    application text                     not null,
    provider    text                     not null,
    host        text                     not null,
    version     text                     not null,
    state       text                     not null,
    metadata    jsonb,
    created_at  timestamp with time zone not null,
    updated_at  timestamp with time zone not null
);

create unique index idx_a_p_h on updates (application, provider, host);

commit;
