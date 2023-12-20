package ssr

import (
	"net/http"

	"site/internal/controller/http/response"
	"site/pkg/logger"
)

type handler struct {
	render ServerSideRender
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log := logger.Load(r.Context())
	resp, err := h.render.Render(r.Context(), r.URL, r.Cookies())
	if err != nil {
		log.Error(err, "can't render")
		response.InternalServerError(rw, r)
		return
	}
	rw.Header().Set("Content-Type", "text/html; charset=UTF-8")
	if _, err := rw.Write(resp); err != nil {
		log.Error(err, "can't write response")
		response.InternalServerError(rw, r)
	}
}

func NewHandler(render ServerSideRender) *handler {
	return &handler{
		render: render,
	}
}
