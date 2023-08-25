package handler

import (
	"github.com/go-logr/logr"
	"net/http"
)

type Interface interface {
	ExtractData(ctx *Context) (interface{}, error)
	Validate(ctx *Context, data interface{}) (*Validation, error)
	Process(ctx *Context, data interface{}) (interface{}, error)
}

type Handler struct {
	Interface
	logger logr.Logger
	opts   handlerOptions
}

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(h.logger, w, r)
	data, err := h.Interface.ExtractData(ctx)
	if nil != err {
		h.logger.Error(err, "can't extract data")
		ctx.RespondWithError(http.StatusBadRequest)
		return
	}
	validation, err := h.Interface.Validate(ctx, data)
	if nil != err {
		h.logger.Error(err, "can't validate data")
		ctx.RespondWithError(http.StatusInternalServerError)
		return
	}
	if !validation.IsValid {
		h.logger.Info("data not valid", "errors", validation.Errors)
		ctx.RespondWithError(http.StatusBadRequest, validation.Errors...)
		return
	}
	response, err := h.Interface.Process(ctx, data)
	if nil != err {
		h.logger.Error(err, "can't process data")
		ctx.RespondWithError(http.StatusInternalServerError)
		return
	}
	if nil != response && h.opts.serializeResultToJSON {
		ctx.respond(http.StatusOK, response)
	}
}

func Create(logger logr.Logger, handler Interface, opts ...Option) http.HandlerFunc {
	h := &Handler{
		Interface: handler,
		logger:    logger,
		opts:      defaultHandlerOptions(),
	}
	for _, opt := range opts {
		opt.apply(&h.opts)
	}
	return h.Handler
}
