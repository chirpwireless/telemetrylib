package telemetrylib

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"os"
)

func InitLogger() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	// Add span context attributes when Context is passed to logging calls.
	instrumentedHandler := handlerWithSpanContext(jsonHandler)
	// Set this handler as the global slog handler.
	slog.SetDefault(slog.New(instrumentedHandler))
}

func handlerWithSpanContext(handler slog.Handler) *spanContextLogHandler {
	return &spanContextLogHandler{Handler: handler}
}

// spanContextLogHandler is an slog.Handler which adds attributes from the
// span context.
type spanContextLogHandler struct {
	slog.Handler
}

// Handle overrides slog.Handler's Handle method. This adds attributes from the
// span context to the slog.Record.
func (t *spanContextLogHandler) Handle(ctx context.Context, record slog.Record) error {
	// Get the SpanContext from the golang Context.
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		// Add trace context attributes following Cloud Logging structured log format described
		// in https://cloud.google.com/logging/docs/structured-logging#special-payload-fields
		record.AddAttrs(
			slog.Any("logging.googleapis.com/trace", spanContext.TraceID()),
		)
		record.AddAttrs(
			slog.Any("logging.googleapis.com/spanId", spanContext.SpanID()),
		)
		record.AddAttrs(
			slog.Bool("logging.googleapis.com/trace_sampled", spanContext.TraceFlags().IsSampled()),
		)
	}
	return t.Handler.Handle(ctx, record)
}
