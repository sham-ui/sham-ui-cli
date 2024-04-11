package metrics

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const RouteName = "metrics"

func Setup(router *mux.Router) {
	router.
		Name(RouteName).
		Methods("GET").
		Path("/metrics").
		Handler(promhttp.Handler())
}
