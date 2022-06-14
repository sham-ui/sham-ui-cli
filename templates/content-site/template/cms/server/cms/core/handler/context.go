package handler

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"cms/core/sessions"
)

// Context container for handler data (Request, ResponseWriter, session)
type Context struct {
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
		log.Errorf("encode json message fail: %s", err)
	}
}

func newContext(w http.ResponseWriter, r *http.Request, sessionsStore *sessions.Store) *Context {
	return &Context{
		Request:       r,
		Response:      w,
		sessionsStore: sessionsStore,
	}
}
