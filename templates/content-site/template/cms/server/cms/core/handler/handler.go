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
	opts   handlerOptions
	logger logr.Logger
}

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(h.logger, w, r, h.opts.sessionStore)
	if h.opts.onlyForAuthenticated {
		session, err := ctx.GetSession()
		if nil != err {
			h.logger.Error(err, "can't get session")
			ctx.RespondWithError(http.StatusInternalServerError)
			return
		}
		if nil == session {
			ctx.RespondWithError(http.StatusUnauthorized, "Session Expired. Log out and log back in.")
			return
		}
		if h.opts.onlyForSuperuser && !session.IsSuperuser {
			ctx.RespondWithError(http.StatusForbidden, "Allowed only for superuser")
			return
		}
	}
	data, err := h.Interface.ExtractData(ctx)
	if nil != err {
		h.logger.V(1).Error(err, "can't extract data")
		ctx.RespondWithError(http.StatusBadRequest)
		return
	}
	validation, err := h.Interface.Validate(ctx, data)
	if nil != err {
		h.logger.V(1).Error(err, "can't validate data")
		ctx.RespondWithError(http.StatusInternalServerError)
		return
	}
	if !validation.IsValid {
		h.logger.V(2).Info("data not valid", "errors", validation.Errors)
		ctx.RespondWithError(http.StatusBadRequest, validation.Errors...)
		return
	}
	response, err := h.Interface.Process(ctx, data)
	if nil != err {
		h.logger.Error(err, "can't process request")
		ctx.RespondWithError(http.StatusInternalServerError)
		return
	}
	if h.opts.serializeResultToJSON {
		ctx.respond(http.StatusOK, response)
	}
}

func Create(logger logr.Logger, handler Interface, opts ...Option) http.HandlerFunc {
	h := &Handler{
		Interface: handler,
		opts:      defaultHandlerOptions(),
		logger:    logger,
	}
	for _, opt := range opts {
		opt.apply(&h.opts)
	}
	return h.Handler
}
