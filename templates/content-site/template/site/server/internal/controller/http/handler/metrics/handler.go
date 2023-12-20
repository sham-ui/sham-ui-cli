package metrics

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Setup(router *mux.Router) {
	router.
		Methods("GET").
		Path("/metrics").
		Handler(promhttp.Handler())
}
