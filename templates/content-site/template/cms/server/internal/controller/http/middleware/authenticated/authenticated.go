package authenticated

import (
	"cms/internal/controller/http/request"
	"cms/internal/controller/http/response"
	"cms/internal/model"
	"cms/pkg/logger"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, ok := request.SessionFromContext(r.Context())
		if !ok {
			logger.Load(r.Context()).Error(model.ErrSessionNotAuthenticated, "session not authenticated")
			response.WithError(rw, r, http.StatusForbidden, "not authenticated")
			return
		}
		next.ServeHTTP(rw, r)
	})
}
