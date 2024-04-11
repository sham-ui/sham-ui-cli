package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type spanExporter struct {
	sdktrace.SpanExporter
}

func (se *spanExporter) String() string {
	return "span-exporter"
}

func (se *spanExporter) GracefulShutdown(ctx context.Context) error {
	if err := se.SpanExporter.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown span exporter: %w", err)
	}
	return nil
}

func NewExporter(cfg Config) (*spanExporter, error) {
	var exporter sdktrace.SpanExporter
	var err error
	if cfg.Endpoint == "" {
		exporter, err = stdouttrace.New()
	} else {
		exporter, err = otlptracehttp.New(
			context.Background(),
			otlptracehttp.WithInsecure(), // use http & not https
			otlptracehttp.WithEndpoint(cfg.Endpoint),
			otlptracehttp.WithURLPath(cfg.Path),
			otlptracehttp.WithHeaders(map[string]string{
				"Authorization": "Basic " + cfg.Authorization,
			}),
		)
	}
	if err != nil {
		return nil, fmt.Errorf("create span exporter: %w", err)
	}
	return &spanExporter{SpanExporter: exporter}, nil
}
