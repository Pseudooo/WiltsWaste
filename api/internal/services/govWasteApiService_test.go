package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Pseudooo/WiltsWaste/internal/db"
	"github.com/Pseudooo/WiltsWaste/internal/models"
	"github.com/Pseudooo/WiltsWaste/internal/services"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGovWasteApiServiceQueries struct {
	mock.Mock
}

func (m *mockGovWasteApiServiceQueries) GetAddressListCacheRow(ctx context.Context, postcode string) ([]byte, error) {
	args := m.Called(ctx, postcode)

	byteArgs := args.Get(0)
	if byteArgs == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).([]byte), args.Error(1)
	}
}

func (m *mockGovWasteApiServiceQueries) InsertAddressListCacheRow(ctx context.Context, params db.InsertAddressListCacheRowParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *mockGovWasteApiServiceQueries) GetWasteCollectionEventsCacheRow(ctx context.Context, params db.GetWasteCollectionEventsCacheRowParams) ([]byte, error) {
	args := m.Called(ctx, params)

	bytesResponse := args.Get(0)
	if bytesResponse == nil {
		return nil, args.Error(1)
	} else {
		return bytesResponse.([]byte), args.Error(1)
	}
}

func (m *mockGovWasteApiServiceQueries) InsertWasteCollectionEventsCacheRow(ctx context.Context, params db.InsertWasteCollectionEventsCacheRowParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

type mockGovWasteApiRepository struct {
	mock.Mock
}

func (m *mockGovWasteApiRepository) GetAddressListFromPostcode(ctx context.Context, postcode string) ([]models.PostcodeAddress, error) {
	args := m.Called(ctx, postcode)
	addressListResult := args.Get(0)
	if addressListResult == nil {
		return nil, args.Error(1)
	} else {
		return addressListResult.([]models.PostcodeAddress), args.Error(1)
	}
}

func (m *mockGovWasteApiRepository) GetWasteCollectionEvents(ctx context.Context, uprn string, postcode string, month int, year int) ([]models.WasteCollectionEvent, error) {
	args := m.Called(ctx, uprn, postcode, month, year)

	wasteCollectionEvents := args.Get(0)
	if wasteCollectionEvents == nil {
		return nil, args.Error(1)
	} else {
		return wasteCollectionEvents.([]models.WasteCollectionEvent), args.Error(1)
	}
}

func TestGetAddressList_GivenCacheEmpty_WhenCalled_ThenCallsApiAndCachesResult(t *testing.T) {
	// Arrange
	addressListResultsToReturn := []models.PostcodeAddress{
		{
			Address: "Number 1, Far Far Away",
			Uprn:    "1",
		},
		{
			Address: "Number 2, Far Far Away",
			Uprn:    "2",
		},
		{
			Address: "Number 3, Far Far Away",
			Uprn:    "3",
		},
	}

	mockQueries := new(mockGovWasteApiServiceQueries)
	mockQueries.On("GetAddressListCacheRow", mock.Anything, "FFA").
		Return(nil, pgx.ErrNoRows)
	mockQueries.On("InsertAddressListCacheRow", mock.Anything, mock.MatchedBy(
		func(params db.InsertAddressListCacheRowParams) bool {
			return len(params.AddressList) > 0
		})).
		Return(nil)

	mockRepository := new(mockGovWasteApiRepository)
	mockRepository.On("GetAddressListFromPostcode", mock.Anything, "FFA").
		Return(addressListResultsToReturn, nil)

	service := services.NewGovWasteApiService(mockRepository, mockQueries)

	// Act
	addressListResults, err := service.GetAddressListFromPostcode(context.Background(), "FFA")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, addressListResults, 3)
	assert.Equal(t, addressListResults[0].Address, "Number 1, Far Far Away")
	assert.Equal(t, addressListResults[0].Uprn, "1")
	assert.Equal(t, addressListResults[1].Address, "Number 2, Far Far Away")
	assert.Equal(t, addressListResults[1].Uprn, "2")
	assert.Equal(t, addressListResults[2].Address, "Number 3, Far Far Away")
	assert.Equal(t, addressListResults[2].Uprn, "3")

	mockQueries.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestGetAddressList_GivenCachePopulated_WhenCalled_ThenReturnsDataFromCache_AndApiNotCalled(t *testing.T) {
	// Arrange
	addressListCachedJson := []byte(`[
		{ "address": "Number 1, Far Far Away", "uprn": "1" },
		{ "address": "Number 2, Far Far Away", "uprn": "2" },
		{ "address": "Number 3, Far Far Away", "uprn": "3" }
	]`)

	mockQueries := new(mockGovWasteApiServiceQueries)
	mockQueries.On("GetAddressListCacheRow", mock.Anything, "FFA").
		Return(addressListCachedJson, nil)

	mockRepository := new(mockGovWasteApiRepository)

	service := services.NewGovWasteApiService(mockRepository, mockQueries)

	// Act
	addressListResults, err := service.GetAddressListFromPostcode(context.Background(), "FFA")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, addressListResults, 3)
	assert.Equal(t, addressListResults[0].Address, "Number 1, Far Far Away")
	assert.Equal(t, addressListResults[0].Uprn, "1")
	assert.Equal(t, addressListResults[1].Address, "Number 2, Far Far Away")
	assert.Equal(t, addressListResults[1].Uprn, "2")
	assert.Equal(t, addressListResults[2].Address, "Number 3, Far Far Away")
	assert.Equal(t, addressListResults[2].Uprn, "3")

	mockQueries.AssertExpectations(t)
	mockRepository.AssertNotCalled(t, "GetAddressListFromPostcode", mock.Anything, mock.Anything)
	mockQueries.AssertNotCalled(t, "InsertAddressListCacheRow", mock.Anything, mock.Anything)
}

func TestGetAddressList_GivenCachePopulatedWithInvalidJson_WhenCalled_ThenCachedDataDiscardedAndFallsBackToApi_AndCachesResponse(t *testing.T) {
	// Arrange
	addressListJsonCache := []byte(`[
		{ "foo": "bar" };
	]`)
	addressListResultsToReturn := []models.PostcodeAddress{
		{
			Address: "Number 1, Far Far Away",
			Uprn:    "1",
		},
		{
			Address: "Number 2, Far Far Away",
			Uprn:    "2",
		},
		{
			Address: "Number 3, Far Far Away",
			Uprn:    "3",
		},
	}

	mockQueries := new(mockGovWasteApiServiceQueries)
	mockQueries.On("GetAddressListCacheRow", mock.Anything, "FFA").
		Return(addressListJsonCache, nil)
	mockQueries.On("InsertAddressListCacheRow", mock.Anything, mock.Anything).
		Return(nil)

	mockRepository := new(mockGovWasteApiRepository)
	mockRepository.On("GetAddressListFromPostcode", mock.Anything, "FFA").
		Return(addressListResultsToReturn, nil)

	service := services.NewGovWasteApiService(mockRepository, mockQueries)

	// Act
	addressListResults, err := service.GetAddressListFromPostcode(context.Background(), "FFA")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, addressListResults, 3)
	assert.Equal(t, addressListResults[0].Address, "Number 1, Far Far Away")
	assert.Equal(t, addressListResults[0].Uprn, "1")
	assert.Equal(t, addressListResults[1].Address, "Number 2, Far Far Away")
	assert.Equal(t, addressListResults[1].Uprn, "2")
	assert.Equal(t, addressListResults[2].Address, "Number 3, Far Far Away")
	assert.Equal(t, addressListResults[2].Uprn, "3")

	mockQueries.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestGetAddressList_GivenCacheLookupFails_WhenCalled_ThenFallsBackToApiAndReturnsResponse_AndCachesResponse(t *testing.T) {
	// Arrange
	addressListResultsToReturn := []models.PostcodeAddress{
		{
			Address: "Number 1, Far Far Away",
			Uprn:    "1",
		},
		{
			Address: "Number 2, Far Far Away",
			Uprn:    "2",
		},
		{
			Address: "Number 3, Far Far Away",
			Uprn:    "3",
		},
	}

	mockQueries := new(mockGovWasteApiServiceQueries)
	mockQueries.On("GetAddressListCacheRow", mock.Anything, "FFA").
		Return(nil, errors.New("Something went OOPSIE when querying the db"))
	mockQueries.On("InsertAddressListCacheRow", mock.Anything, mock.Anything).
		Return(nil)

	mockRepository := new(mockGovWasteApiRepository)
	mockRepository.On("GetAddressListFromPostcode", mock.Anything, "FFA").
		Return(addressListResultsToReturn, nil)

	service := services.NewGovWasteApiService(mockRepository, mockQueries)

	// Act
	addressListResults, err := service.GetAddressListFromPostcode(context.Background(), "FFA")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, addressListResults, 3)
	assert.Equal(t, addressListResults[0].Address, "Number 1, Far Far Away")
	assert.Equal(t, addressListResults[0].Uprn, "1")
	assert.Equal(t, addressListResults[1].Address, "Number 2, Far Far Away")
	assert.Equal(t, addressListResults[1].Uprn, "2")
	assert.Equal(t, addressListResults[2].Address, "Number 3, Far Far Away")
	assert.Equal(t, addressListResults[2].Uprn, "3")

	mockQueries.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestGetAddressList_GivenCacheEmpty_WhenCalledAndCacheInsertFails_ThenDataStillReturned(t *testing.T) {
	// Arrange
	addressListResultsToReturn := []models.PostcodeAddress{
		{
			Address: "Number 1, Far Far Away",
			Uprn:    "1",
		},
		{
			Address: "Number 2, Far Far Away",
			Uprn:    "2",
		},
		{
			Address: "Number 3, Far Far Away",
			Uprn:    "3",
		},
	}

	insertError := errors.New("Something went OOPSIE writing to the db")
	mockQueries := new(mockGovWasteApiServiceQueries)
	mockQueries.On("GetAddressListCacheRow", mock.Anything, "FFA").
		Return(nil, pgx.ErrNoRows)
	mockQueries.On("InsertAddressListCacheRow", mock.Anything, mock.Anything).
		Return(insertError)

	mockRepository := new(mockGovWasteApiRepository)
	mockRepository.On("GetAddressListFromPostcode", mock.Anything, "FFA").
		Return(addressListResultsToReturn, nil)

	service := services.NewGovWasteApiService(mockRepository, mockQueries)

	// Act
	addressListResults, err := service.GetAddressListFromPostcode(context.Background(), "FFA")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, addressListResults, 3)
	assert.Equal(t, addressListResults[0].Address, "Number 1, Far Far Away")
	assert.Equal(t, addressListResults[0].Uprn, "1")
	assert.Equal(t, addressListResults[1].Address, "Number 2, Far Far Away")
	assert.Equal(t, addressListResults[1].Uprn, "2")
	assert.Equal(t, addressListResults[2].Address, "Number 3, Far Far Away")
	assert.Equal(t, addressListResults[2].Uprn, "3")

	mockQueries.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestGetAddressList_GivenCacheEmpty_WhenCalledAndApiFails_ThenErrorReturned_AndInsertNotAttempted(t *testing.T) {
	// Arrange
	repositoryError := errors.New("Something went OOPSIE calling the repository")
	mockQueries := new(mockGovWasteApiServiceQueries)
	mockQueries.On("GetAddressListCacheRow", mock.Anything, "FFA").
		Return(nil, pgx.ErrNoRows)

	mockRepository := new(mockGovWasteApiRepository)
	mockRepository.On("GetAddressListFromPostcode", mock.Anything, "FFA").
		Return(nil, repositoryError)

	service := services.NewGovWasteApiService(mockRepository, mockQueries)

	// Act
	addressListResults, err := service.GetAddressListFromPostcode(context.Background(), "FFA")

	// Assert
	assert.Nil(t, addressListResults)
	assert.Error(t, err)

	mockQueries.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
	mockQueries.AssertNotCalled(t, "InsertAddressListCacheRow", mock.Anything, mock.Anything)
}

func TestGetWasteCollectionEvents_GivenCacheEmpty_WhenCallsApi_ThenReturnsResultAndInsertsToCache(t *testing.T) {
	// Arrange
	wasteCollectionEventsToReturn := []models.WasteCollectionEvent{
		{
			Title: "Normal Rubbish",
			Date:  time.Now().AddDate(0, 0, 1),
		},
		{
			Title: "Recycling",
			Date:  time.Now().AddDate(0, 0, 5),
		},
	}

	mockQueries := new(mockGovWasteApiServiceQueries)
	mockQueries.On("GetWasteCollectionEventsCacheRow", mock.Anything, mock.Anything).
		Return(nil, pgx.ErrNoRows)
	mockQueries.On("InsertWasteCollectionEventsCacheRow", mock.Anything, mock.Anything).
		Return(nil)

	mockRepository := new(mockGovWasteApiRepository)
	mockRepository.On("GetWasteCollectionEvents", mock.Anything, "1", "FFA", 01, 2026).
		Return(wasteCollectionEventsToReturn, nil)

	service := services.NewGovWasteApiService(mockRepository, mockQueries)

	// Act
	wasteCollectionEvents, err := service.GetWasteCollectionEvents(context.Background(), "1", "FFA", 01, 2026)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, wasteCollectionEvents, 2)
	assert.Equal(t, wasteCollectionEvents[0].Title, "Normal Rubbish")
	assert.Equal(t, wasteCollectionEvents[1].Title, "Recycling")

	mockQueries.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestGetWasteCollectionEvents_GivenCachePopulated_WhenCalled_ThenDoesntCallApiReturnsCachedResponse(t *testing.T) {
	// Arrange
	cachedWasteCollectionEventsJson := []byte(`[
		{ "Title": "Normal Rubbish", "Date": "2026-11-10T00:00:00Z" },
		{ "Title": "Recycling", "Date": "2026-11-15T00:00:00Z" }
	]`)

	mockQueries := new(mockGovWasteApiServiceQueries)
	mockQueries.On("GetWasteCollectionEventsCacheRow", mock.Anything, mock.Anything).
		Return(cachedWasteCollectionEventsJson, nil)

	mockRepository := new(mockGovWasteApiRepository)

	service := services.NewGovWasteApiService(mockRepository, mockQueries)

	// Act
	wasteCollectionEvents, err := service.GetWasteCollectionEvents(context.Background(), "1", "FFA", 01, 2026)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, wasteCollectionEvents, 2)
	assert.Equal(t, wasteCollectionEvents[0].Title, "Normal Rubbish")
	assert.Equal(t, wasteCollectionEvents[1].Title, "Recycling")

	mockQueries.AssertExpectations(t)
	mockRepository.AssertNotCalled(t, "GetWasteCollectionEvents", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	mockQueries.AssertNotCalled(t, "InsertWasteCollectionEventsCacheRow", mock.Anything, mock.Anything)
}

func TestGetWasteCollectionEvents_GivenCacheCorrupted_WhenCalled_ThenFallsBackToApi_AndInsertsNewCacheRow(t *testing.T) {
	// Arrange
	corruptedCacheJson := []byte(`[
		{ "foo:" "bar";}
	]`)
	wasteCollectionEventsToReturn := []models.WasteCollectionEvent{
		{
			Title: "Normal Rubbish",
			Date:  time.Now().AddDate(0, 0, 1),
		},
		{
			Title: "Recycling",
			Date:  time.Now().AddDate(0, 0, 5),
		},
	}

	mockQueries := new(mockGovWasteApiServiceQueries)
	mockQueries.On("GetWasteCollectionEventsCacheRow", mock.Anything, mock.Anything).
		Return(corruptedCacheJson, nil)
	mockQueries.On("InsertWasteCollectionEventsCacheRow", mock.Anything, mock.Anything).
		Return(nil)

	mockRepository := new(mockGovWasteApiRepository)
	mockRepository.On("GetWasteCollectionEvents", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(wasteCollectionEventsToReturn, nil)

	service := services.NewGovWasteApiService(mockRepository, mockQueries)

	// Act
	wasteCollectionEvents, err := service.GetWasteCollectionEvents(context.Background(), "1", "FFA", 01, 2026)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, wasteCollectionEvents, 2)
	assert.Equal(t, wasteCollectionEvents[0].Title, "Normal Rubbish")
	assert.Equal(t, wasteCollectionEvents[1].Title, "Recycling")

	mockQueries.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestGetWasteCollectionEvents_GivenCacheLookupFails_WhenCalled_ThenFallsBackToApi_AndInsertsNewCacheRow(t *testing.T) {
	// Arrange
	wasteCollectionEventsToReturn := []models.WasteCollectionEvent{
		{
			Title: "Normal Rubbish",
			Date:  time.Now().AddDate(0, 0, 1),
		},
		{
			Title: "Recycling",
			Date:  time.Now().AddDate(0, 0, 5),
		},
	}

	mockQueries := new(mockGovWasteApiServiceQueries)
	mockQueries.On("GetWasteCollectionEventsCacheRow", mock.Anything, mock.Anything).
		Return(nil, errors.New("Something went OOPSIE querying the db"))
	mockQueries.On("InsertWasteCollectionEventsCacheRow", mock.Anything, mock.Anything).
		Return(nil)

	mockRepository := new(mockGovWasteApiRepository)
	mockRepository.On("GetWasteCollectionEvents", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(wasteCollectionEventsToReturn, nil)

	service := services.NewGovWasteApiService(mockRepository, mockQueries)

	// Act
	wasteCollectionEvents, err := service.GetWasteCollectionEvents(context.Background(), "1", "FFA", 01, 2026)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, wasteCollectionEvents, 2)
	assert.Equal(t, wasteCollectionEvents[0].Title, "Normal Rubbish")
	assert.Equal(t, wasteCollectionEvents[1].Title, "Recycling")

	mockQueries.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestGetWasteCollectionEvents_GivenCacheInsertFails_WhenCalled_ThenReturnsResponse(t *testing.T) {
	// Arrange
	wasteCollectionEventsToReturn := []models.WasteCollectionEvent{
		{
			Title: "Normal Rubbish",
			Date:  time.Now().AddDate(0, 0, 1),
		},
		{
			Title: "Recycling",
			Date:  time.Now().AddDate(0, 0, 5),
		},
	}

	mockQueries := new(mockGovWasteApiServiceQueries)
	mockQueries.On("GetWasteCollectionEventsCacheRow", mock.Anything, mock.Anything).
		Return(nil, pgx.ErrNoRows)
	mockQueries.On("InsertWasteCollectionEventsCacheRow", mock.Anything, mock.Anything).
		Return(errors.New("Something went OOPSIE writing to the db"))

	mockRepository := new(mockGovWasteApiRepository)
	mockRepository.On("GetWasteCollectionEvents", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(wasteCollectionEventsToReturn, nil)

	service := services.NewGovWasteApiService(mockRepository, mockQueries)

	// Act
	wasteCollectionEvents, err := service.GetWasteCollectionEvents(context.Background(), "1", "FFA", 01, 2026)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, wasteCollectionEvents, 2)
	assert.Equal(t, wasteCollectionEvents[0].Title, "Normal Rubbish")
	assert.Equal(t, wasteCollectionEvents[1].Title, "Recycling")

	mockQueries.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestGetWasteCollectionEvents_GivenApiReturnsError_WhenCalled_ThenReturnsErrror(t *testing.T) {
	// Arrange
	mockQueries := new(mockGovWasteApiServiceQueries)
	mockQueries.On("GetWasteCollectionEventsCacheRow", mock.Anything, mock.Anything).
		Return(nil, pgx.ErrNoRows)

	mockRepository := new(mockGovWasteApiRepository)
	mockRepository.On("GetWasteCollectionEvents", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("Something is bad in the response"))

	service := services.NewGovWasteApiService(mockRepository, mockQueries)

	// Act
	wasteCollectionEvents, err := service.GetWasteCollectionEvents(context.Background(), "1", "FFA", 01, 2026)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, wasteCollectionEvents)

	mockQueries.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
	mockQueries.AssertNotCalled(t, "InsertWasteCollectionEventsCacheRow", mock.Anything, mock.Anything)
}
