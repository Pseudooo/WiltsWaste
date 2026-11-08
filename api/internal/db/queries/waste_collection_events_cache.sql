-- name: InsertWasteCollectionEventsCacheRow :exec
insert into waste_collection_events_cache (uprn, postcode, month, year, waste_collection_events, created_at, expires_at)
    values (@uprn, @postcode, @month, @year, @waste_collection_events, NOW(), @expires_at);

-- name: GetWasteCollectionEventsCacheRow :one
select waste_collection_events
    from waste_collection_events_cache
    where uprn = @uprn 
        and postcode = @postcode 
        and year = @year 
        and month = @month
        and expires_at > now()
    order by created_at desc
    limit 1;
