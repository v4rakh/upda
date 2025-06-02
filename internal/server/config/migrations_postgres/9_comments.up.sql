create table if not exists comments
(
    id          uuid                     not null constraint uni_comments_id primary key,
    author      text                     not null,
    content     text                     not null,
    update_id    uuid                    not null constraint fk_comments_update references updates on update cascade on delete cascade,
    created_at  timestamp with time zone not null,
    updated_at  timestamp with time zone not null
);

commit;
