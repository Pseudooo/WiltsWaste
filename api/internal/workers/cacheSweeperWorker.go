package workers

import (
	"context"
	"log/slog"
	"time"

	"github.com/Pseudooo/WiltsWaste/internal/telemetry"
)

type CacheSweeperWorker interface {
	StartBackgroundWorker(ctx context.Context)
}

type CacheSweeperWorkerQueries interface {
	DeleteExpiredAddressListCacheRows(ctx context.Context) (int64, error)
	DeleteExpiredWasteCollectionEventsCacheRows(ctx context.Context) (int64, error)
}

type cacheSweeperWorker struct {
	queries CacheSweeperWorkerQueries
}

func NewCacheSweeperWorker(queries CacheSweeperWorkerQueries) CacheSweeperWorker {
	return &cacheSweeperWorker{
		queries: queries,
	}
}

func (w *cacheSweeperWorker) StartBackgroundWorker(ctx context.Context) {
	telemetry.Logger.Debug("Cache Sweeper Worker starting...")

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			telemetry.Logger.Info("Background worker terminating...")
			return

		case <-ticker.C:
			w.deleteExpiredCacheRows(ctx)
			return
		}
	}
}

func (w *cacheSweeperWorker) deleteExpiredCacheRows(ctx context.Context) {
	logger := telemetry.Logger

	addressListCacheRowsCleaned, err := w.queries.DeleteExpiredAddressListCacheRows(ctx)
	if err != nil {
		logger.Error("Failed to run delete query for AddressList cache table", slog.Any("error", err))
	}

	wasteCollectionEventsCacheRowsCleaned, err := w.queries.DeleteExpiredWasteCollectionEventsCacheRows(ctx)
	if err != nil {
		logger.Error("Failed to run delete query for WasteCollectionEvents cache table", slog.Any("error", err))
	}

	if addressListCacheRowsCleaned == 0 && wasteCollectionEventsCacheRowsCleaned == 0 {
		logger.Debug("No cache rows to delete")
		return
	}

	logger.Info(
		"Cache Tables Cleared",
		slog.Int64("address_list_rows_cleaned", addressListCacheRowsCleaned),
		slog.Int64("waste_collection_rows_cleaned", wasteCollectionEventsCacheRowsCleaned),
	)
}
