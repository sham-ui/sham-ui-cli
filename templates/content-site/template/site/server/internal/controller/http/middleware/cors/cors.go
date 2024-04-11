package cors

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Setup(router *mux.Router) {
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			rw.Header().Set("Access-Control-Allow-Credentials", "true")
			next.ServeHTTP(rw, r)
		})
	})
}
