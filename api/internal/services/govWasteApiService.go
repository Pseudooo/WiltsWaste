package services

import (
	"context"

	"github.com/Pseudooo/WiltsWaste/internal/models"
	govWasteApi "github.com/Pseudooo/WiltsWaste/internal/repositories"
)

type GovWasteApiService interface {
	GetAddressListFromPostcode(ctx context.Context, postcode string) ([]models.PostcodeAddress, error)
}

type govWasteApiService struct {
	repository *govWasteApi.GovWasteApiRepository
}

func NewGovWasteApiService(govWasteApiRepository *govWasteApi.GovWasteApiRepository) GovWasteApiService {
	return &govWasteApiService{
		repository: govWasteApiRepository,
	}
}

func (s *govWasteApiService) GetAddressListFromPostcode(ctx context.Context, postcode string) ([]models.PostcodeAddress, error) {
	return s.repository.GetAddressListFromPostcode(ctx, postcode)
}
