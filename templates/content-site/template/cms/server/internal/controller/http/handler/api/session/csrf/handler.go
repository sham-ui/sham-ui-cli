package csrf

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

const RouteName = "api.session.csrf"

type handler struct {
	CSRFRequestHeader string
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set(h.CSRFRequestHeader, csrf.Token(r))
}

func newHandler(csrfRequestHeader string) *handler {
	return &handler{
		CSRFRequestHeader: csrfRequestHeader,
	}
}

func Setup(router *mux.Router, csrfRequestHeader string) {
	router.
		Name(RouteName).
		Methods(http.MethodGet).
		Path("/csrftoken").
		Handler(newHandler(csrfRequestHeader))
}
