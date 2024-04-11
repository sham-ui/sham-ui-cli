package metrics

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const RouteName = "metrics"

func Setup(router *mux.Router) {
	router.
		Name(RouteName).
		Methods(http.MethodGet).
		Path("/metrics").
		Handler(promhttp.Handler())
}
