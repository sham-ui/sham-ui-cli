package cors

import (
	"github.com/gorilla/mux"
)

func Setup(r *mux.Router) {
	r.Use(mux.CORSMethodMiddleware(r))
}
