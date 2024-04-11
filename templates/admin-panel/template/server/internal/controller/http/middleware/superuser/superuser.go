package superuser

import (
	"net/http"
	"{{ shortName }}/internal/controller/http/request"
	"{{ shortName }}/internal/controller/http/response"
	"{{ shortName }}/internal/model"
	"{{ shortName }}/pkg/logger"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		sess, ok := request.SessionFromContext(r.Context())
		if !ok || !sess.IsSuperuser {
			logger.Load(r.Context()).Error(model.ErrSessionNotSuperuser, "not superuser")
			response.WithError(rw, r, http.StatusForbidden, "not allowed")
			return
		}
		next.ServeHTTP(rw, r)
	})
}
