package handler

import (
	"cms/core/sessions"
	"encoding/json"
	"github.com/go-logr/logr"
	"net/http"
)

// Context container for handler data (Request, ResponseWriter, session)
type Context struct {
	logger           logr.Logger
	sessionsStore    *sessions.Store
	hasCachedSession bool
	session          *sessions.Session
	Request          *http.Request
	Response         http.ResponseWriter
}

type responseError struct {
	Status   string   `json:"Status"`
	Messages []string `json:"Messages,omitempty"`
}

// GetSession return pointer to session (if exists)
func (ctx *Context) GetSession() (*sessions.Session, error) {
	var err error
	if !ctx.hasCachedSession {
		ctx.session, err = ctx.sessionsStore.GetSession(ctx.Request)
		ctx.hasCachedSession = true
	}
	return ctx.session, err
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
		ctx.logger.Error(err, "encode json message fail")
	}
}

func newContext(logger logr.Logger, w http.ResponseWriter, r *http.Request, sessionsStore *sessions.Store) *Context {
	return &Context{
		logger:        logger,
		Request:       r,
		Response:      w,
		sessionsStore: sessionsStore,
	}
}
