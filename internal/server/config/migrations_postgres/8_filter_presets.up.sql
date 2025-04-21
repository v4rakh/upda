create table if not exists filter_presets (
    id uuid not null constraint uni_filter_presets_id primary key,
    "type" text not null,
    label text not null,
    parameters text not null,
    color text,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null
);
commit;
