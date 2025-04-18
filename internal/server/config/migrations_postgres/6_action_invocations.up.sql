create table if not exists action_invocations
(
    id          uuid                     not null
        constraint uni_action_invocations_id primary key,
    retry_count bigint default 1         not null,
    state       text                     not null,
    message     text,
    event_id    uuid                     not null
        constraint fk_action_invocations_event references events on update cascade on delete cascade,
    action_id   uuid                     not null
        constraint fk_action_invocations_action references actions on update cascade on delete cascade,
    created_at  timestamp with time zone not null,
    updated_at  timestamp with time zone not null
);

commit;
