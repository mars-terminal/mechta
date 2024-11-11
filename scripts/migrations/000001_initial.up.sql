create table links (
    id uuid,
    short_link text not null,
    target_url text not null,
    expire_at timestamptz not null,
    last_access timestamptz,
    access_count int default 0,
    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    deleted_at timestamptz,

    primary key (id)
);

create unique index on links (short_link);
create index on links (target_url);
create index on links (created_at desc);

