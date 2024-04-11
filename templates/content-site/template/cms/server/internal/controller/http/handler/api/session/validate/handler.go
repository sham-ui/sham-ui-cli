package validate

import (
	"cms/internal/controller/http/request"
	"cms/internal/controller/http/response"
	"net/http"

	"github.com/gorilla/mux"
)

const RouteName = "api.session.validate"

type validateResponse struct {
	Name        string `json:"Name"`
	Email       string `json:"Email"`
	IsSuperuser bool   `json:"IsSuperuser"`
}

func handler(rw http.ResponseWriter, r *http.Request) {
	sess, ok := request.SessionFromContext(r.Context())
	if !ok {
		response.WithError(rw, r, http.StatusUnauthorized, "not authenticated")
		return
	}
	response.JSON(rw, r, http.StatusOK, validateResponse{
		Name:        sess.Name,
		Email:       sess.Email,
		IsSuperuser: sess.IsSuperuser,
	})
}

func Setup(router *mux.Router) {
	router.
		Name(RouteName).
		Methods(http.MethodGet).
		Path("/validsession").
		HandlerFunc(handler)
}
