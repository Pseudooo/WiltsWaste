package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"uuid"

	"github.com/Pseudooo/WiltsWaste/internal/addresslist"
	"github.com/Pseudooo/WiltsWaste/internal/calendar"
	"github.com/Pseudooo/WiltsWaste/internal/configuration"
	"github.com/Pseudooo/WiltsWaste/internal/db"
	govWasteApi "github.com/Pseudooo/WiltsWaste/internal/repositories"
	"github.com/Pseudooo/WiltsWaste/internal/services"
	"github.com/Pseudooo/WiltsWaste/internal/telemetry"
	externalHttpClient "github.com/Pseudooo/WiltsWaste/internal/util"
	"github.com/Pseudooo/WiltsWaste/internal/workers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const LoggerKey = "slog_logger"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownMetrics, err := telemetry.InitMetrics(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := shutdownMetrics(ctx); err != nil {
			fmt.Printf("Error shutting down metrics%v\n", err)
		}
	}()

	globalLogger := telemetry.Logger

	databaseConfiguration, err := configuration.ReadDatabaseConfigurationFromEnv()
	if err != nil {
		panic(fmt.Errorf("Failed to read database configuration: %w", err))
	}

	sqlDb, err := databaseConfiguration.CreateDatabase(ctx)
	if err != nil {
		panic(err)
	}
	defer sqlDb.Close()

	queries := db.New(sqlDb)

	cacheSweeper := workers.NewCacheSweeperWorker(queries)
	go cacheSweeper.StartBackgroundWorker(ctx)

	govWasteApiClient := externalHttpClient.CreateExternalHttpClient("gov-waste-api")
	govWasteApiRepository := govWasteApi.CreateRepository(govWasteApiClient, "https://ilambassadorformsprod.azurewebsites.net")

	govWasteApiService := services.NewGovWasteApiService(govWasteApiRepository, queries)

	calendarHandler := &calendar.CalendarHandler{
		Logger:          globalLogger,
		WasteApiService: govWasteApiService,
	}
	addressListHandler := &addresslist.AddressListHandler{
		Logger:          globalLogger,
		WasteApiService: govWasteApiService,
	}

	router := gin.New()
	router.Use(SlogMiddleware(globalLogger))
	router.Use(otelgin.Middleware("waste-collection-api"))
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	router.GET("/calendar", calendarHandler.GetCalendar)
	router.GET("/addresslist", addressListHandler.GetAddressList)

	router.Run(":8080")
}

func SlogMiddleware(baseLogger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		correlationId := c.GetHeader("X-Correlation-Id")
		if correlationId == "" {
			correlationId = uuid.New().String()
		}
		c.Header("X-Correlation-Id", correlationId)

		requestLogger := baseLogger.With(
			slog.String("correlation_id", correlationId),
		)
		c.Set(LoggerKey, requestLogger)

		c.Next()

		elapsed := time.Since(start)
		response_code := strconv.Itoa(c.Writer.Status())

		requestLogger.Info(
			"Request Completed",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("query", c.Request.URL.RawQuery),
			slog.Duration("latency", elapsed),
			slog.String("status_code", response_code),
		)

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				requestLogger.Error("error", slog.String("error", err.Error()))
			}
		}
	}
}
