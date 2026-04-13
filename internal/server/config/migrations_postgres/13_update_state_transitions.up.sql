create table if not exists update_state_transitions
(
    id            uuid                     not null
        constraint uni_update_state_transitions_id primary key,
    from_state_id uuid                     not null
        references update_state_definitions on delete cascade,
    to_state_id   uuid                     not null
        references update_state_definitions on delete cascade,
    created_at    timestamp with time zone not null,
    updated_at    timestamp with time zone not null,
    constraint uni_update_state_transitions_from_to unique (from_state_id, to_state_id)
);

-- Seed all transitions between default states
insert into update_state_transitions (id, from_state_id, to_state_id, created_at, updated_at)
select gen_random_uuid(), f.id, t.id, now(), now()
from update_state_definitions f, update_state_definitions t
where f.name != t.name;

commit;
