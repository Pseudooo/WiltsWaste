-- name: InsertAddressListCacheRow :exec
insert into address_list_cache (postcode, created_at, expires_at, address_list)
    values (@postcode, NOW(), @expires_at, @address_list);

-- name: GetAddressListCacheRow :one
select address_list
    from address_list_cache
    where postcode = @postcode and expires_at < now()
    order by created_at desc;
