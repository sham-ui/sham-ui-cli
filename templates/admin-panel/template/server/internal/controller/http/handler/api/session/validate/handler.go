package validate

import (
	"net/http"
	"{{ shortName }}/internal/controller/http/request"
	"{{ shortName }}/internal/controller/http/response"

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
