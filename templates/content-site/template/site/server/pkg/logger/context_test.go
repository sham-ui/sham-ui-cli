package logger

import (
	"context"
	"testing"

	"site/pkg/asserts"

	"github.com/go-logr/logr/testr"
)

func TestContextLogger(t *testing.T) {
	logger := testr.New(t)
	ctx := Save(context.Background(), logger)
	asserts.Equals(t, logger, Load(ctx))
}
