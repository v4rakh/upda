create table if not exists update_state_definitions
(
    id          uuid                     not null
        constraint uni_update_state_definitions_id primary key,
    name        text                     not null
        constraint uni_update_state_definitions_name unique,
    label       text                     not null,
    color       text                     not null default 'gray',
    icon        text                     not null default 'QuestionOutlined',
    description text,
    is_initial  boolean                  not null default false,
    skip_on_new_version  boolean         not null default false,
    sort_order  integer                  not null default 0,
    created_at  timestamp with time zone not null,
    updated_at  timestamp with time zone not null
);

-- Seed default states
insert into update_state_definitions (id, name, label, color, icon, description, is_initial, skip_on_new_version, sort_order, created_at, updated_at)
values
    (gen_random_uuid(), 'pending', 'Pending', '#00ccff', 'InteractionOutlined', 'Update is awaiting action', true, false,0, now(), now()),
    (gen_random_uuid(), 'approved', 'Approved', 'limegreen', 'CheckCircleOutlined', 'Update has been handled successfully', false, false,1, now(), now()),
    (gen_random_uuid(), 'ignored', 'Ignored', 'goldenrod', 'StopOutlined', 'Future changes to this update are being ignored', false,true, 2, now(), now());

commit;
