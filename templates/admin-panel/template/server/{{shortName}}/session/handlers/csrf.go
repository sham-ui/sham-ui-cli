package handlers

import (
	"github.com/go-logr/logr"
	"github.com/gorilla/csrf"
	"net/http"
	"{{shortName}}/core/handler"
)

// csrfToken will generate a CSRF Token
func csrfToken(ctx *handler.Context, _ interface{}) (interface{}, error) {
	ctx.Response.Header().Set("X-CSRF-Token", csrf.Token(ctx.Request))
	return nil, nil
}

func NewCsrfTokenHandler(logger logr.Logger) http.HandlerFunc {
	return handler.CreateFromProcessFunc(logger, csrfToken, handler.WithoutSerializeResultToJSON())
}
