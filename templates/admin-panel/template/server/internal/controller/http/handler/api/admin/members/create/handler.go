package create

import (
	"errors"
	"net/http"
	"{{ shortName }}/internal/controller/http/request"
	"{{ shortName }}/internal/controller/http/response"
	"{{ shortName }}/internal/model"
	"{{ shortName }}/pkg/logger"
	"{{ shortName }}/pkg/validation"
	"strings"

	"github.com/gorilla/mux"
)

const RouteName = "api.admin.members.create"

type requestData struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsSuperUser bool   `json:"is_superuser"`
}

type responseData struct {
	Status string `json:"Status"`
}

type handler struct {
	memberService   MemberService
	passwordService PasswordService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
	if len(data.Password) == 0 {
		valid.AddErrors("Password must not be empty.")
	}
	if valid.HasErrors() {
		response.BadRequest(rw, r, valid.Errors...)
		return
	}

	hashedPassword, err := h.passwordService.Hash(ctx, data.Password)
	if err != nil {
		logger.Load(ctx).Error(err, "failed to hash password")
		response.InternalServerError(rw, r)
		return
	}

	err = h.memberService.Create(ctx, model.MemberWithPassword{
		Member: model.Member{ //nolint:exhaustruct
			Name:        data.Name,
			Email:       data.Email,
			IsSuperuser: data.IsSuperUser,
		},
		HashedPassword: hashedPassword,
	})
	switch {
	case errors.Is(err, model.ErrMemberEmailAlreadyExists):
		response.BadRequest(rw, r, "Email is already in use.")
		return
	case err != nil:
		logger.Load(ctx).Error(err, "failed to create member")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusOK, &responseData{
		Status: "Member created",
	})
}

func newHandler(memberService MemberService, passwordService PasswordService) *handler {
	return &handler{
		memberService:   memberService,
		passwordService: passwordService,
	}
}

func Setup(router *mux.Router, memberService MemberService, passwordService PasswordService) {
	router.
		Name(RouteName).
		Methods(http.MethodPost).
		Handler(newHandler(memberService, passwordService))
}
