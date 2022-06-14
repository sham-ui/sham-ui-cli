package handlers

import (
	"fmt"
	"net/http"
	"cms/core/handler"
	"cms/core/sessions"
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

func NewLogoutHandler(sessionsStore *sessions.Store) http.HandlerFunc {
	return handler.CreateFromProcessFunc(
		logoutHandler,
		handler.WithOnlyForAuthenticated(sessionsStore),
		handler.WithoutSerializeResultToJSON(),
	)
}
