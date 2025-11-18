create table if not exists rastochnoy_write (
    id uuid primary key default gen_random_uuid(),
    key varchar(30) unique not null,
    offsett float unique not null,
    value BOOL NOT NULL
);