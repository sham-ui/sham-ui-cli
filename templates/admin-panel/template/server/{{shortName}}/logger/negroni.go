package logger

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/urfave/negroni"
)

type negroniALogger struct {
	logger logr.Logger
}

func (n negroniALogger) Println(v ...interface{}) {
	n.logger.Info(fmt.Sprint(v...))
}

func (n negroniALogger) Printf(format string, v ...interface{}) {
	n.logger.Info(fmt.Sprintf(format, v...))
}

func CreateNegroniLogger(l logr.Logger) *negroni.Logger {
	l = l.WithName("negroni")
	adapter := &negroniALogger{logger: l}
	loggerMiddleware := negroni.NewLogger()
	loggerMiddleware.ALogger = adapter
	loggerMiddleware.SetFormat("\{{.Status}} | \{{.Duration}} | \{{.Hostname}} | \{{.Method}} \{{.Path}}")
	return loggerMiddleware
}
