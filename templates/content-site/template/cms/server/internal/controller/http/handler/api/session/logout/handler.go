package logout

import (
	"cms/internal/controller/http/middleware/authenticated"
	"cms/internal/controller/http/response"
	"cms/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

const RouteName = "api.session.logout"

type handler struct {
	service Service
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(rw, r); err != nil {
		logger.Load(r.Context()).Error(err, "failed to delete session")
		response.InternalServerError(rw, r)
	}
}

func newHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func Setup(router *mux.Router, service Service) {
	router = router.NewRoute().Subrouter()
	router.Use(
		authenticated.Middleware,
	)
	router.
		Name(RouteName).
		Methods(http.MethodPost).
		Path("/logout").
		Handler(newHandler(service))
}
