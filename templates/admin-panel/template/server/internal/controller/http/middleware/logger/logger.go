package logger

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"go.opentelemetry.io/otel/trace"

	"github.com/go-logr/logr"
)

type logger struct {
	logger logr.Logger
}

func (cl *logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logRw := newLogResponseWriter(rw)

		next.ServeHTTP(logRw, r)

		fields := []any{
			"start", start.Format(time.RFC3339Nano),
			"method", r.Method,
			"hostname", r.Host,
			"request_uri", r.URL.RequestURI(),
			"user_agent", r.UserAgent(),
			"remote_addr", r.RemoteAddr,
			"status", logRw.StatusCode(),
			"duration", time.Since(start),
		}

		routeStr := ""
		route := mux.CurrentRoute(r)
		if route != nil {
			var err error
			routeStr, err = route.GetPathTemplate()
			if err != nil {
				routeStr, err = route.GetPathRegexp()
				if err != nil {
					routeStr = ""
				}
			}
		}
		if routeStr != "" {
			fields = append(fields, "route", routeStr)
		}

		sc := trace.SpanContextFromContext(r.Context())
		if sc.IsValid() {
			fields = append(fields,
				"trace_id", sc.TraceID().String(),
				"span_id", sc.SpanID().String(),
			)
		}

		cl.logger.Info("http request", fields...)
	})
}

func New(l logr.Logger) *logger {
	return &logger{logger: l}
}
