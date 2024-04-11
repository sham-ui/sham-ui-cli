package superuser

import (
	"cms/internal/controller/http/request"
	"cms/internal/controller/http/response"
	"cms/internal/model"
	"cms/pkg/logger"
	"net/http"
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
