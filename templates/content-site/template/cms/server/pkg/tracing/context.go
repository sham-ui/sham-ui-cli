package tracing

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func SpanFromContext(ctx context.Context) trace.Span {
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		return span
	}
	return noop.Span{Span: nil}
}

func ProviderFromContext(ctx context.Context) trace.TracerProvider {
	return SpanFromContext(ctx).TracerProvider()
}

func TracerFromContext(ctx context.Context, scopeName string, opts ...trace.TracerOption) trace.Tracer {
	return ProviderFromContext(ctx).Tracer(scopeName, opts...)
}

func StartSpan(ctx context.Context, scopeName, operation string, opts ...trace.SpanStartOption) (context.Context, trace.Span) { //nolint:lll
	return TracerFromContext(ctx, scopeName).Start(ctx, scopeName+"."+operation, opts...)
}
