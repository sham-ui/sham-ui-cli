package handlers

import (
	"fmt"
	"github.com/go-logr/logr"
	"net/http"
	"{{shortName}}/core/handler"
	"{{shortName}}/core/sessions"
)

func logoutHandler(ctx *handler.Context, _ interface{}) (interface{}, error) {
	session, _ := ctx.GetSession()
	rawSession := session.GetRawSession()
	// Revoke users authentication
	rawSession.Values["authenticated"] = false
	rawSession.Options.MaxAge = -1
	err := rawSession.Save(ctx.Request, ctx.Response)
	if nil != err {
		return nil, fmt.Errorf("can't save session: %s", err)
	}
	return nil, nil
}

func NewLogoutHandler(logger logr.Logger, sessionsStore *sessions.Store) http.HandlerFunc {
	return handler.CreateFromProcessFunc(
		logger,
		logoutHandler,
		handler.WithOnlyForAuthenticated(sessionsStore),
		handler.WithoutSerializeResultToJSON(),
	)
}
