create table if not exists sessions
(
    id         text not null
    primary key,
    data       text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    expires_at timestamp with time zone
);

commit;
