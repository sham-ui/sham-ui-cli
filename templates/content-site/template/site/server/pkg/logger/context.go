package logger

import (
	"context"

	"github.com/go-logr/logr"
)

func Save(ctx context.Context, logger logr.Logger) context.Context {
	return logr.NewContext(ctx, logger)
}

func Load(ctx context.Context) logr.Logger {
	return logr.FromContextOrDiscard(ctx)
}
