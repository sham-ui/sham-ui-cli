package handler

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Context container for handler data (Request, ResponseWriter)
type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
}

type responseError struct {
	Status   string   `json:"Status"`
	Messages []string `json:"Messages,omitempty"`
}

// respondWithError send error to client
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
		log.Errorf("encode json message fail: %s", err)
	}
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Request:  r,
		Response: w,
	}
}
