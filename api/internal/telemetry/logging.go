package telemetry

import (
	"context"
	"log/slog"
	"os"
)

var Logger = createGlobalLogger()

func createGlobalLogger() *slog.Logger {
	globalLogAttributes := []slog.Attr{
		slog.String("service.name", "wilts-waste-api"),
		slog.String("service.version", "0.0.0"),
	}
	loggingHandler := slog.NewJSONHandler(os.Stdout, nil).WithAttrs(globalLogAttributes)
	return slog.New(loggingHandler)
}

func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	if ctxLogger := ctx.Value("slog_logger"); ctxLogger != nil {
		if logger, ok := ctxLogger.(*slog.Logger); ok {
			return logger
		}
	}

	return Logger
}
