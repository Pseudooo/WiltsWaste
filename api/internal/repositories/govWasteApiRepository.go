package govWasteApi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Pseudooo/WiltsWaste/internal/models"
	externalHttpClient "github.com/Pseudooo/WiltsWaste/internal/util"
	"github.com/PuerkitoBio/goquery"
)

type GovApiError struct {
	StatusCode   int
	Url          string
	Uprn         string
	Postcode     string
	Month        int
	Year         int
	ResponseBody string
}

func (e *GovApiError) Error() string {
	return fmt.Errorf("Api request failed with status %d", e.StatusCode).Error()
}

func (e *GovApiError) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("status_code", e.StatusCode),
		slog.String("url", e.Url),
		slog.String("uprn", e.Uprn),
		slog.String("postcode", e.Postcode),
		slog.Int("month", e.Month),
		slog.Int("year", e.Year),
		slog.String("response_body", e.ResponseBody),
	)
}

type GovWasteApiRepository struct {
	httpClient *externalHttpClient.ExternalHttpClient
	baseUrl    string
}

func CreateRepository(client *externalHttpClient.ExternalHttpClient, baseUrl string) *GovWasteApiRepository {
	return &GovWasteApiRepository{
		httpClient: client,
		baseUrl:    baseUrl,
	}
}

type GovApiAddressListError struct {
	StatusCode   int
	Url          string
	Postcode     string
	ResponseBody string
}

func (e *GovApiAddressListError) Error() string {
	return fmt.Errorf("Gov Api AddressList request failed with status: %d", e.StatusCode).Error()
}

func (e *GovApiAddressListError) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("status_code", e.StatusCode),
		slog.String("url", e.Url),
		slog.String("postcode", e.Postcode),
		slog.String("response_body", e.ResponseBody),
	)
}

type govApiAddressListResponse struct {
	Model     govApiAddressListModel `json:"Model"`
	IsSuccess bool                   `json:"IsSuccess"`
	Message   string                 `json:"Message"`
}

type govApiAddressListModel struct {
	PostcodeAddresses []govApiPostCodeAddress `json:"PostcodeAddresses"`
}

type govApiPostCodeAddress struct {
	Address string `json:"Address"`
	Uprn    string `json:"UPRN"`
}

func (r *GovWasteApiRepository) GetAddressListFromPostcode(ctx context.Context, postcode string) ([]models.PostcodeAddress, error) {
	endpoint := r.baseUrl + "/wastecollectiondays/addresslist"

	data := url.Values{}
	data.Set("Postcode", postcode)

	res, err := r.httpClient.Do(
		ctx,
		"get-address-list",
		http.MethodPost,
		endpoint,
		strings.NewReader(data.Encode()),
		externalHttpClient.WithHeader("User-Agent", "github/Pseudooo/WiltsWaste"),
		externalHttpClient.WithHeader("Content-Type", "application/x-www-form-urlencoded"),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(res.Body, 1024)
		bytes, _ := io.ReadAll(limitedReader)

		return nil, &GovApiAddressListError{
			StatusCode:   res.StatusCode,
			Url:          endpoint,
			Postcode:     postcode,
			ResponseBody: string(bytes),
		}
	}

	var result govApiAddressListResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var mappedResult = make([]models.PostcodeAddress, len(result.Model.PostcodeAddresses))
	for _, v := range result.Model.PostcodeAddresses {
		mappedResult = append(mappedResult, models.PostcodeAddress{
			Address: v.Address,
			Uprn:    v.Uprn,
		})
	}

	return mappedResult, nil
}

func (r *GovWasteApiRepository) GetWasteCollectionEvents(ctx context.Context, uprn string, postcode string, month int, year int) ([]models.WasteCollectionEvent, error) {
	endpoint := r.baseUrl + "/wastecollectiondays/wastecollectioncalendar"

	data := url.Values{}
	data.Set("Postcode", postcode)
	data.Set("Month", strconv.Itoa(month))
	data.Set("Year", strconv.Itoa(year))
	data.Set("Uprn", uprn)

	res, err := r.httpClient.Do(
		ctx,
		"get-waste-events",
		http.MethodPost,
		endpoint,
		strings.NewReader(data.Encode()),
		externalHttpClient.WithHeader("User-Agent", "github/Pseudooo/WasteCollectionCalendar"),
		externalHttpClient.WithHeader("Content-Type", "application/x-www-form-urlencoded"),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(res.Body, 1024)
		bytes, _ := io.ReadAll(limitedReader)

		return nil, &GovApiError{
			StatusCode:   res.StatusCode,
			Url:          res.Request.URL.RawPath,
			Uprn:         uprn,
			Postcode:     postcode,
			Month:        month,
			Year:         year,
			ResponseBody: string(bytes),
		}
	}

	return parseWasteCollectionEventsFromReader(res.Body)
}

func parseWasteCollectionEventsFromReader(reader io.Reader) ([]models.WasteCollectionEvent, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var events []models.WasteCollectionEvent

	doc.Find(".cal-month-box .rc-event-container").Each(func(i int, s *goquery.Selection) {
		eventTitle := strings.TrimSpace(s.Find("span").Text())

		calInnerParent := s.Closest(".cal-inner")

		dateStr, exists := calInnerParent.Find(".day-no").Attr("data-cal-date")
		parsedTime, err := time.Parse("2006-01-02T15:04:05", dateStr)

		if exists && eventTitle != "" && err == nil {
			events = append(events, models.WasteCollectionEvent{
				Title: eventTitle,
				Date:  parsedTime,
			})
		}
	})

	return events, nil
}
