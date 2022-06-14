package client

import (
	"encoding/json"
	"log"
	"net/http/httptest"
)

type ResponseWrapper struct {
	Response *httptest.ResponseRecorder
}

func (r *ResponseWrapper) JSON() map[string]interface{} {
	data := map[string]interface{}{}
	err := json.Unmarshal(r.Response.Body.Bytes(), &data)
	if nil != err {
		log.Fatalf("unmarshal: %s", err)
	}
	return data
}

func (r *ResponseWrapper) Text() string {
	return r.Response.Body.String()
}

func newResponseWrapper(resp *httptest.ResponseRecorder) *ResponseWrapper {
	return &ResponseWrapper{Response: resp}
}
