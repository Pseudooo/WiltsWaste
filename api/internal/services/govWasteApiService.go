package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Pseudooo/WiltsWaste/internal/db"
	"github.com/Pseudooo/WiltsWaste/internal/models"
	govWasteApi "github.com/Pseudooo/WiltsWaste/internal/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type GovWasteApiService interface {
	GetAddressListFromPostcode(ctx context.Context, postcode string) ([]models.PostcodeAddress, error)
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
	cachedResponse, err := s.queries.GetAddressListCacheRow(ctx, postcode)
	if err == nil {
		var addressList []models.PostcodeAddress
		if err := json.Unmarshal(cachedResponse, &addressList); err != nil {
			return addressList, nil
		}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		// TODO: Log db failure here
	}

	addressList, err := s.repository.GetAddressListFromPostcode(ctx, postcode)
	if err != nil {
		return nil, err
	}

	addressListJsonBytes, err := json.Marshal(addressList)
	if err != nil {
		err = s.queries.InsertAddressListCacheRow(ctx, db.InsertAddressListCacheRowParams{
			Postcode:    postcode,
			ExpiresAt:   pgtype.Timestamp{Valid: true, Time: time.Now().AddDate(0, 0, 1)},
			AddressList: addressListJsonBytes,
		})
		if err != nil {
			// TODO: More logging
		}
	}

	return addressList, nil
}
