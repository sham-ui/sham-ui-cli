package cors

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	commonAllowedHeaders = "Cookie, Content-Type, Content-Length, Accept-Encoding, Accept, Origin, Cache-Control"
	commonExposedHeaders = "Set-Cookie"
)

func Setup(csrfRequestHeader string, router *mux.Router) {
	allowedHeaders := commonAllowedHeaders + ", " + csrfRequestHeader
	exposedHeaders := commonExposedHeaders + ", " + csrfRequestHeader

	// Registry OPTIONS handler
	router.Methods(http.MethodOptions).HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
	})

	// Enable CORS
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			headers := rw.Header()
			headers.Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			headers.Set("Access-Control-Allow-Headers", allowedHeaders)
			headers.Set("Access-Control-Allow-Credentials", "true")
			if r.Method != http.MethodOptions {
				headers.Set("Access-Control-Expose-Headers", exposedHeaders)
				headers.Set("Access-Control-Allow-Methods", r.Method)
			}
			next.ServeHTTP(rw, r)
		})
	})
}
