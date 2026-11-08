-- +goose Up
create table if not exists address_list_cache(
    id uuid primary key default uuidv7(),
    postcode text not null,
    created_at timestamp,
    expires_at timestamp,
    address_list json
);

-- +goose Down
drop table if exists address_list_cache;