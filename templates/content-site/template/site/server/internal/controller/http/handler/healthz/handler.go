package healthz

import (
	"net/http"

	"github.com/gorilla/mux"
)

func handler(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
}

func Setup(router *mux.Router) {
	router.
		Methods("GET").
		Path("/healthz").
		HandlerFunc(handler)
}
