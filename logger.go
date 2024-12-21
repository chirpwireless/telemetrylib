package telemetrylib

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"os"
)

func InitLogger(options *slog.HandlerOptions) {
	if options == nil {
		options = &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, options)
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
		// Add trace context attributes following OpenTelemetry structured log format described
		// in https://opentelemetry.io/docs/concepts/signals/traces/
		record.AddAttrs(
			slog.Any("trace_id", spanContext.TraceID()),
		)
		record.AddAttrs(
			slog.Any("span_id", spanContext.SpanID()),
		)
		record.AddAttrs(
			slog.Bool("trace_sampled", spanContext.TraceFlags().IsSampled()),
		)
	}
	return t.Handler.Handle(ctx, record)
}
