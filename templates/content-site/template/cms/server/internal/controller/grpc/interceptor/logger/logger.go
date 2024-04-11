package logger

import (
	"cms/pkg/tracing"
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

func interceptorLogger(log logr.Logger) logging.LoggerFunc {
	return func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		fields = tracing.AppendTraceId(ctx, fields)
		log := log.WithValues(fields...)
		switch lvl {
		case logging.LevelDebug:
			log.V(3).Info(msg) //nolint:gomnd
		case logging.LevelInfo:
			log.V(2).Info(msg) //nolint:gomnd
		case logging.LevelWarn:
			log.V(1).Info(msg)
		case logging.LevelError:
			log.V(0).Info(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	}
}

func New(l logr.Logger) grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(interceptorLogger(l))
}
