package context_logger

import (
	"net/http"

	"go.opentelemetry.io/otel/trace"

	"site/pkg/logger"

	"github.com/go-logr/logr"
)

type contextLogger struct {
	logger logr.Logger
}

func (cl *contextLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fields := []any{
			"method", r.Method,
			"path", r.URL.Path,
			"host", r.Host,
			"request_uri", r.URL.RequestURI(),
			"user_agent", r.UserAgent(),
			"remote_addr", r.RemoteAddr,
		}
		sc := trace.SpanContextFromContext(r.Context())
		if sc.IsValid() {
			fields = append(fields,
				"trace_id", sc.TraceID().String(),
				"span_id", sc.SpanID().String(),
			)
		}
		requestLogger := cl.logger.WithName("request").WithValues(fields...)
		ctx := logger.Save(r.Context(), requestLogger)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

func New(logger logr.Logger) *contextLogger {
	return &contextLogger{logger: logger}
}
