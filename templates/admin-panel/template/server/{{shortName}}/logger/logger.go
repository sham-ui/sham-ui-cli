package logger

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
	"os"
	"time"
)

func zerologUnixTimestampHook(e *zerolog.Event, level zerolog.Level, message string) {
	e.Int64("timestamp", time.Now().UnixNano())
}

func NewLogger(level int) logr.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	zerologr.NameFieldName = "logger"
	zerologr.NameSeparator = "/"
	zerologr.SetMaxV(level)

	zl := zerolog.New(os.Stdout).
		With().
		Caller().
		Timestamp().
		Logger().
		Hook(zerolog.HookFunc(zerologUnixTimestampHook))
	return zerologr.New(&zl)
}
