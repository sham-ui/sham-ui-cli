package handlers

import (
	"github.com/gorilla/csrf"
	"net/http"
	"cms/core/handler"
)

// csrfToken will generate a CSRF Token
func csrfToken(ctx *handler.Context, _ interface{}) (interface{}, error) {
	ctx.Response.Header().Set("X-CSRF-Token", csrf.Token(ctx.Request))
	return nil, nil
}

func NewCsrfTokenHandler() http.HandlerFunc {
	return handler.CreateFromProcessFunc(csrfToken, handler.WithoutSerializeResultToJSON())
}
