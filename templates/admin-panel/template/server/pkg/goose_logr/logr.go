package goose_logr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/pressly/goose/v3"
)

type logrAdapter logr.Logger

//nolint:exhaustruct
var _ goose.Logger = logrAdapter{}

func (l logrAdapter) Fatalf(format string, v ...interface{}) {
	str := strings.TrimSpace(fmt.Sprintf(format, v...))
	err := errors.New(str) //nolint:goerr113
	logr.Logger(l).Error(err, "goose fatal error")
}

func (l logrAdapter) Printf(format string, v ...interface{}) {
	str := strings.TrimSpace(fmt.Sprintf(format, v...))
	logr.Logger(l).Info(str)
}

func NewLogger(log logr.Logger) logrAdapter {
	return logrAdapter(log)
}
