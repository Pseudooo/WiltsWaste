-- +goose Up
create table if not exists address_list_cache(
    id uuid primary key default uuidv7(),
    postcode text not null,
    created_at timestamp,
    expires_at timestamp,
    address_list json
);

create table if not exists waste_collection_events_cache(
    id uuid primary key default uuidv7(),
    uprn text not null,
    postcode text not null,
    month int not null,
    year int not null,
    created_at timestamp,
    expires_at timestamp,
    waste_collection_events json
)

-- +goose Down
drop table if exists address_list_cache;