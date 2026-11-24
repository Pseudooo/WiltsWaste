package calendar

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Pseudooo/WiltsWaste/internal/models"
	"github.com/Pseudooo/WiltsWaste/internal/services"
	"github.com/Pseudooo/WiltsWaste/internal/telemetry"
	ics "github.com/arran4/golang-ical"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type CalendarHandler struct {
	Logger          *slog.Logger
	WasteApiService services.GovWasteApiService
}

type CalendarQuery struct {
	Uprn     string `form:"uprn" binding:"required"`
	Postcode string `form:"postcode" binding:"required"`
}

func (h *CalendarHandler) GetCalendar(ctx *gin.Context) {
	logger := telemetry.GetLoggerFromContext(ctx)

	var query CalendarQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wasteCollectionEvents, err := h.getWasteCollectionEvents(ctx, query.Uprn, query.Postcode)
	if err != nil {
		logger.Error("Failed to lookup waste collection events", slog.Any("error", err))
		ctx.AbortWithStatus(500)
		return
	}

	ctx.Header("Vary", "Accept")

	accept := ctx.GetHeader("Accept")
	if strings.Contains(accept, "application/json") {
		// Client is requesting json data so events can be returned directly
		ctx.Header("Content-Type", "application/json")
		ctx.JSON(http.StatusOK, wasteCollectionEvents)
		return
	}

	// Otherwise default to the text/calendar ICS content
	calendar := buildCalendarFromWasteCollectionEvents(wasteCollectionEvents)

	ctx.Header("Content-Type", "text/calendar; charset=utf-8")
	ctx.Header("Content-Disposition", `attachment; filename="calendar.ics"`)
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	ctx.String(http.StatusOK, calendar.Serialize(ics.WithNewLineWindows))
}

func (h *CalendarHandler) getWasteCollectionEvents(ctx context.Context, uprn string, postcode string) ([]models.WasteCollectionEvent, error) {
	currentTime := time.Now()
	results := make([][]models.WasteCollectionEvent, 3)

	group, groupCtx := errgroup.WithContext(ctx)

	for i := range 3 {
		index := i
		group.Go(func() error {
			evalTime := currentTime.AddDate(0, index, 0)
			events, err := h.WasteApiService.GetWasteCollectionEvents(groupCtx, uprn, postcode, int(evalTime.Month()), evalTime.Year())
			if err != nil {
				return err
			}

			results[index] = events
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	var allEvents []models.WasteCollectionEvent
	for _, events := range results {
		allEvents = append(allEvents, events...)
	}

	return allEvents, nil
}
