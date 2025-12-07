package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

type contextKey string

const TraceIDKey contextKey = "trace_id"

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		r.AddAttrs(slog.String("trace_id", traceID))
	}
	return h.Handler.Handle(ctx, r)
}

// NewLogger creates a new slog.Logger with the specified level.
// level can be "DEBUG", "INFO", "WARN", "ERROR".
// Defaults to INFO if an invalid or empty level is provided.
// If filePath is provided, logs will be written to that file with rotation.
// If filePath is empty, logs will be written to stdout.
func NewLogger(level string, filePath string) (*slog.Logger, error) {
	var logLevel slog.Level
	switch strings.ToUpper(level) {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	var writer io.Writer

	if filePath != "" {
		writer = &lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		}
	} else {
		writer = os.Stdout
	}

	baseHandler := slog.NewTextHandler(writer, &slog.HandlerOptions{
		Level: logLevel,
	})

	handler := &ContextHandler{Handler: baseHandler}

	return slog.New(handler), nil
}
