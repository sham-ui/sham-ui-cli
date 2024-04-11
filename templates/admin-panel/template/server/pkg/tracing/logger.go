package tracing

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func AppendTraceId(ctx context.Context, fields []any) []any {
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		fields = append(fields,
			"trace_id", sc.TraceID().String(),
			"span_id", sc.SpanID().String(),
		)
	}
	return fields
}
