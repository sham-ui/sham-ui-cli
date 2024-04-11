package healthz

import (
	"net/http"

	"github.com/gorilla/mux"
)

const RouteName = "healthz"

func handler(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
}

func Setup(router *mux.Router) {
	router.
		Name(RouteName).
		Methods(http.MethodGet).
		Path("/healthz").
		HandlerFunc(handler)
}
