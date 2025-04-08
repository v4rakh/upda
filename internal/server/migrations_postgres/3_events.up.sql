create table if not exists events
(
    id         uuid                     not null
        constraint uni_events_id primary key,
    name       text                     not null,
    state      text                     not null,
    payload    jsonb,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null
);


commit;
