package addresslist

import (
	"log/slog"
	"net/http"

	"github.com/Pseudooo/WiltsWaste/internal/services"
	"github.com/gin-gonic/gin"
)

type AddressListHandler struct {
	Logger          *slog.Logger
	WasteApiService services.GovWasteApiService
}

type AddressListQuery struct {
	Postcode string `form:"postcode" binding:"required"`
}

func (h *AddressListHandler) GetAddressList(ctx *gin.Context) {
	logger := h.Logger
	if ctxLogger, exists := ctx.Get("slog_logger"); exists {
		if l, ok := ctxLogger.(*slog.Logger); ok {
			logger = l
		}
	}

	var query AddressListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.WasteApiService.GetAddressListFromPostcode(ctx, query.Postcode)
	if err != nil {
		logger.Error("An error occurred when calling gov address list endpoint", slog.Any("error", err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, results)
}
