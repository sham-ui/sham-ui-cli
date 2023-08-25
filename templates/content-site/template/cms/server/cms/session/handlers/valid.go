package handlers

import (
	"cms/core/handler"
	"cms/core/sessions"
	"github.com/go-logr/logr"
	"net/http"
)

type validSessionResponse struct {
	Name        string
	Email       string
	IsSuperuser bool
}

// Process checks that the session is valid and can user can make requests
func validSessionHandler(ctx *handler.Context, _ interface{}) (interface{}, error) {
	session, _ := ctx.GetSession()
	return &validSessionResponse{
		Name:        session.Name,
		Email:       session.Email,
		IsSuperuser: session.IsSuperuser,
	}, nil
}

func NewValidSessionHandler(logger logr.Logger, sessionsStore *sessions.Store) http.HandlerFunc {
	return handler.CreateFromProcessFunc(logger, validSessionHandler, handler.WithOnlyForAuthenticated(sessionsStore))
}
