package client

import (
	"net/http/httptest"
)

type ResponseWrapper struct {
	Response *httptest.ResponseRecorder
}

func (r *ResponseWrapper) Text() string {
	return r.Response.Body.String()
}

func newResponseWrapper(resp *httptest.ResponseRecorder) *ResponseWrapper {
	return &ResponseWrapper{Response: resp}
}
