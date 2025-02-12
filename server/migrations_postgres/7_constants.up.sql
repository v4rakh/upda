create table if not exists constants
(
    id         uuid                     not null
        constraint uni_constants_id primary key,
    key        text                     not null
        constraint uni_constants_key unique,
    value      text                     not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null
);


commit;
