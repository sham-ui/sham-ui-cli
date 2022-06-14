package handler

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Interface interface {
	ExtractData(ctx *Context) (interface{}, error)
	Validate(ctx *Context, data interface{}) (*Validation, error)
	Process(ctx *Context, data interface{}) (interface{}, error)
}

type Handler struct {
	Interface
	opts handlerOptions
}

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(w, r, h.opts.sessionStore)
	if h.opts.onlyForAuthenticated {
		session, err := ctx.GetSession()
		if nil != err {
			log.Errorf("can't get session: %s", err)
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
		log.Infof("can't extract data: %s", err)
		ctx.RespondWithError(http.StatusBadRequest)
		return
	}
	validation, err := h.Interface.Validate(ctx, data)
	if nil != err {
		log.Errorf("can't validate data: %s", err)
		ctx.RespondWithError(http.StatusInternalServerError)
		return
	}
	if !validation.IsValid {
		log.WithField("errors", validation.Errors).Info("data not valid")
		ctx.RespondWithError(http.StatusBadRequest, validation.Errors...)
		return
	}
	response, err := h.Interface.Process(ctx, data)
	if nil != err {
		log.Errorf("can't process request: %s", err)
		ctx.RespondWithError(http.StatusInternalServerError)
		return
	}
	if h.opts.serializeResultToJSON {
		ctx.respond(http.StatusOK, response)
	}
}

func Create(handler Interface, opts ...Option) http.HandlerFunc {
	h := &Handler{
		Interface: handler,
		opts:      defaultHandlerOptions(),
	}
	for _, opt := range opts {
		opt.apply(&h.opts)
	}
	return h.Handler
}
