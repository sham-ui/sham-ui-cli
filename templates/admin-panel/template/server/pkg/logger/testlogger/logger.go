package testlogger

import (
	"strings"

	"github.com/go-logr/logr"
)

type logger struct {
	Messages []Message
}

func (l *logger) WithName(
	_ string,
) logr.LogSink {
	return l
}
func (l *logger) WithValues(_ ...any) logr.LogSink {
	return l
}
func (l *logger) Init(_ logr.RuntimeInfo) {}

func (l *logger) Enabled(_ int) bool { return true }
func (l *logger) Info(level int, message string, keyValues ...any) {
	l.Messages = append(l.Messages, newMessage(level, message, keyValues...))
}
func (l *logger) Error(err error, message string, keyValues ...any) {
	l.Messages = append(l.Messages, newMessage(0, message, append(keyValues, "error", err.Error())...))
}

func (l *logger) String() string {
	var out strings.Builder
	last := len(l.Messages) - 1
	for i, m := range l.Messages {
		out.WriteString(m.String())
		if i < last {
			out.WriteString("\n")
		}
	}
	return out.String()
}

func NewLogger() *logger {
	return &logger{
		Messages: make([]Message, 0),
	}
}
