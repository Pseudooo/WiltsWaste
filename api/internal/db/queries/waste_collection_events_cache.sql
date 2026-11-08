-- name: InsertWasteCollectionEventsCacheRow :exec
insert into waste_collection_events_cache (uprn, postcode, month, year, waste_collection_events)
    values (@uprn, @postcode, @month, @year, @waste_collection_events);

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
