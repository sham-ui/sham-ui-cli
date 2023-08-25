package handler

import (
	"encoding/json"
	"github.com/go-logr/logr"
	"net/http"
)

// Context container for handler data (Request, ResponseWriter)
type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	logger   logr.Logger
}

type responseError struct {
	Status   string   `json:"Status"`
	Messages []string `json:"Messages,omitempty"`
}

// RespondWithError send error to client
func (ctx *Context) RespondWithError(statusCode int, messages ...string) {
	msg := &responseError{
		Status:   http.StatusText(statusCode),
		Messages: messages,
	}
	ctx.respond(statusCode, msg)
}

func (ctx *Context) respond(statusCode int, msg interface{}) {
	ctx.Response.WriteHeader(statusCode)
	ctx.Response.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(ctx.Response).Encode(msg)
	if nil != err {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		ctx.logger.Error(err, "encode json message fail")
	}
}

func newContext(logger logr.Logger, w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Request:  r,
		Response: w,
		logger:   logger,
	}
}
