create table if not exists rastochnoy_read (
    id uuid primary key default gen_random_uuid(),
    key varchar(30) unique not null,
    offsett float unique not null,
    value BOOLEAN NOT NULL DEFAULT FALSE
);