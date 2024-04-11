package tracing

import "go.opentelemetry.io/otel/propagation"

func NewPropagator() propagation.TraceContext {
	return propagation.TraceContext{}
}
