package tracing

import (
	"net/http"
	"site/pkg/middleware"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func New(tracerProvider trace.TracerProvider, propagator propagation.TraceContext) func(http.Handler) http.Handler {
	return middleware.Compose(
		otelmux.Middleware("http",
			otelmux.WithTracerProvider(tracerProvider),
			otelmux.WithPropagators(propagator),
			otelmux.WithSpanNameFormatter(func(routeName string, r *http.Request) string {
				if routeName == "/" {
					routeName = r.URL.Path
				}
				if r.URL.RawQuery != "" {
					routeName += "?" + r.URL.RawQuery
				}
				return routeName
			}),
		),
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				defer next.ServeHTTP(rw, r)

				span := trace.SpanFromContext(r.Context())
				if !span.SpanContext().IsValid() {
					return
				}

				vars := mux.Vars(r)
				attrs := make([]attribute.KeyValue, 0, len(vars))
				for k, v := range vars {
					attrs = append(attrs, attribute.String(k, v))
				}

				if len(attrs) > 0 {
					span.SetAttributes(attrs...)
				}
			})
		},
	)
}
