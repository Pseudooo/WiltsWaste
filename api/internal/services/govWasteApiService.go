package services

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/Pseudooo/WiltsWaste/internal/db"
	"github.com/Pseudooo/WiltsWaste/internal/models"
	govWasteApi "github.com/Pseudooo/WiltsWaste/internal/repositories"
	"github.com/Pseudooo/WiltsWaste/internal/telemetry"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type GovWasteApiService interface {
	GetAddressListFromPostcode(ctx context.Context, postcode string) ([]models.PostcodeAddress, error)
	GetWasteCollectionEvents(ctx context.Context, uprn string, postcode string, month int, year int) ([]models.WasteCollectionEvent, error)
}

type govWasteApiService struct {
	repository *govWasteApi.GovWasteApiRepository
	queries    *db.Queries
}

func NewGovWasteApiService(
	govWasteApiRepository *govWasteApi.GovWasteApiRepository,
	queries *db.Queries,
) GovWasteApiService {
	return &govWasteApiService{
		repository: govWasteApiRepository,
		queries:    queries,
	}
}

func (s *govWasteApiService) GetAddressListFromPostcode(ctx context.Context, postcode string) ([]models.PostcodeAddress, error) {
	ctxLogger := telemetry.GetLoggerFromContext(ctx)

	cachedResponse, err := s.queries.GetAddressListCacheRow(ctx, postcode)
	if err == nil {
		var addressList []models.PostcodeAddress
		if err := json.Unmarshal(cachedResponse, &addressList); err != nil {
			return addressList, nil
		}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		ctxLogger.Error("Failed to query address list cache table", slog.Any("error", err))
	}

	addressList, err := s.repository.GetAddressListFromPostcode(ctx, postcode)
	if err != nil {
		return nil, err
	}

	addressListJsonBytes, err := json.Marshal(addressList)
	if err == nil {
		err = s.queries.InsertAddressListCacheRow(ctx, db.InsertAddressListCacheRowParams{
			Postcode:    postcode,
			ExpiresAt:   pgtype.Timestamp{Valid: true, Time: time.Now().AddDate(0, 0, 1)},
			AddressList: addressListJsonBytes,
		})
		if err != nil {
			ctxLogger.Error(
				"Failed to insert address list cache row",
				slog.String("postcode", postcode),
				slog.Any("error", err),
			)
		}
	} else {
		ctxLogger.Error(
			"Failed to serialise address list response",
			slog.Any("error", err),
		)
	}

	return addressList, nil
}

func (s *govWasteApiService) GetWasteCollectionEvents(
	ctx context.Context,
	uprn string,
	postcode string,
	month int,
	year int,
) ([]models.WasteCollectionEvent, error) {
	return s.repository.GetWasteCollectionEvents(ctx, uprn, postcode, month, year)
}
