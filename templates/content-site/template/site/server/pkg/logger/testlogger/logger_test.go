package testlogger

import (
	"site/pkg/asserts"
	"testing"

	"github.com/go-logr/logr"
)

// Compile time check that Logger implements LogSink.
var _ logr.LogSink = NewLogger()

func TestLogger(t *testing.T) {
	// Arrange
	tlog := NewLogger()
	log := logr.New(tlog)

	// Act
	log.Info("info", "key", "value")

	// Assert
	asserts.Equals(t, 1, len(tlog.Messages), "message count")
	asserts.Equals(t, 0, tlog.Messages[0].Level, "message level")
	asserts.Equals(t, "info", tlog.Messages[0].Message, "message")
	asserts.Equals(t, 1, len(tlog.Messages[0].KeyValues), "message key values count")
	asserts.Equals(t, "value", tlog.Messages[0].KeyValues["key"], "message key")
}

func TestLogger_String(t *testing.T) {
	// Arrange
	tlog := NewLogger()
	log := logr.New(tlog)

	// Act
	log.Info("info", "key", "value")

	// Assert
	asserts.JSONEquals(t, `{"level": 0, "message": "info", "key": "value"}`, tlog.String(), "string")
}
