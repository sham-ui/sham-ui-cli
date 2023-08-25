package handler

import (
	"github.com/go-logr/logr"
	"net/http"
)

type ProcessFunc func(ctx *Context, data interface{}) (interface{}, error)

// onlyProcessFuncHandler is simple Interface implementation for wrap ProcessFunc
type onlyProcessFuncHandler struct {
	HandlerWithoutExtractDataAndValidation
	fn ProcessFunc
}

func (h *onlyProcessFuncHandler) Process(ctx *Context, data interface{}) (interface{}, error) {
	return h.fn(ctx, data)
}

func CreateFromProcessFunc(logger logr.Logger, fn ProcessFunc, opts ...Option) http.HandlerFunc {
	return Create(logger, &onlyProcessFuncHandler{
		fn: fn,
	}, opts...)
}
