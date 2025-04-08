create table if not exists actions
(
    id                uuid                     not null
        constraint uni_actions_id primary key,
    label             text                     not null,
    type              text                     not null,
    match_event       text,
    match_application text,
    match_provider    text,
    match_host        text,
    payload           jsonb,
    enabled           boolean                  not null,
    created_at        timestamp with time zone not null,
    updated_at        timestamp with time zone not null
);

commit;
