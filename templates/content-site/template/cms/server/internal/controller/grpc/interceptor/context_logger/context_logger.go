package context_logger

import (
	"cms/pkg/logger"
	"cms/pkg/tracing"
	"context"

	"github.com/go-logr/logr"
	"google.golang.org/grpc"
)

func New(log logr.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		fields := []any{
			"method", info.FullMethod,
		}
		fields = tracing.AppendTraceId(ctx, fields)
		ctx = logger.Save(ctx, log.WithValues(fields...))
		return handler(ctx, req)
	}
}
