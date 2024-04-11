package recovery

import (
	"fmt"
	"net/http"
)

type httpHandlerPanicError struct {
	Err     any
	Request *http.Request
}

func (e *httpHandlerPanicError) Error() string {
	return fmt.Sprintf(
		"http handler panic: %s %s: %v",
		e.Request.Method,
		e.Request.URL.RequestURI(),
		e.Err,
	)
}

func newHttpHandlerPanicError(err any, request *http.Request) error {
	return &httpHandlerPanicError{
		Err:     err,
		Request: request,
	}
}
