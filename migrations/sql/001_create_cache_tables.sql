-- +goose Up
create table if not exists address_list_cache(
    id uuid primary key default uuidv7(),
    postcode text not null,
    created_at timestamp,
    expires_at timestamp,
    address_list json
);

create index idx_address_list_cache_lookup
    on address_list_cache (postcode, created_at desc, expires_at);

create table if not exists waste_collection_events_cache(
    id uuid primary key default uuidv7(),
    uprn text not null,
    postcode text not null,
    month int not null,
    year int not null,
    created_at timestamp,
    expires_at timestamp,
    waste_collection_events json
);

create index idx_waste_collection_events_cache_lookup
    on waste_collection_events_cache (uprn, postcode, year, month, created_at desc, expires_at);

-- +goose Down
drop index if exists idx_address_list_cache_lookup;
drop index if exists idx_waste_collection_events_cache_lookup;
drop table if exists address_list_cache;
drop table if exists waste_collection_events_cache;