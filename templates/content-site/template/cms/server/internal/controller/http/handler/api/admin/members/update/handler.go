package update

import (
	"cms/internal/controller/http/request"
	"cms/internal/controller/http/response"
	"cms/internal/model"
	"cms/pkg/logger"
	"cms/pkg/validation"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const RouteName = "api.admin.members.update"

const idKey = "id"

type requestData struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	IsSuperUser bool   `json:"is_superuser"`
}

type responseData struct {
	Status string `json:"Status"`
}

type handler struct {
	memberService MemberService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := mux.Vars(r)[idKey]
	if !ok {
		response.BadRequest(rw, r, "Empty id")
		return
	}

	data, err := request.DecodeJSON[requestData](r)
	if err != nil {
		response.BadRequest(rw, r, "Invalid JSON")
		return
	}
	data.Name = strings.TrimSpace(data.Name)
	data.Email = strings.TrimSpace(data.Email)

	valid := validation.New()
	if len(data.Name) == 0 {
		valid.AddErrors("Name must not be empty.")
	}
	if len(data.Email) == 0 {
		valid.AddErrors("Email must not be empty.")
	}
	if valid.HasErrors() {
		response.BadRequest(rw, r, valid.Errors...)
		return
	}

	err = h.memberService.Update(ctx, model.Member{
		ID:          model.MemberID(id),
		Name:        data.Name,
		Email:       data.Email,
		IsSuperuser: data.IsSuperUser,
	})
	switch {
	case errors.Is(err, model.ErrMemberEmailAlreadyExists):
		response.BadRequest(rw, r, "Email is already in use.")
		return
	case err != nil:
		logger.Load(ctx).Error(err, "failed to update member")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusOK, &responseData{
		Status: "Member updated",
	})
}

func newHandler(memberService MemberService) *handler {
	return &handler{
		memberService: memberService,
	}
}

func Setup(router *mux.Router, memberService MemberService) {
	router.
		Name(RouteName).
		Methods(http.MethodPut).
		Path("/{" + idKey + ":[0-9]+}").
		Handler(newHandler(memberService))
}
