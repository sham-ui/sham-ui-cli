package signup

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"{{ shortName }}/internal/controller/http/request"
	"{{ shortName }}/internal/controller/http/response"
	"{{ shortName }}/internal/model"
	"{{ shortName }}/pkg/logger"
	"{{ shortName }}/pkg/validation"
	"strings"
)

const RouteName = "api.members.signup"

type requestData struct {
	Name      string `json:"Name"`
	Email     string `json:"Email"`
	Password  string `json:"Password"`
	Password2 string `json:"Password2"`
}

type responseData struct {
	Status string `json:"Status"`
}

type handler struct {
	memberService   MemberService
	passwordService PasswordService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) { //nolint:funlen,cyclop
	ctx := r.Context()

	_, ok := request.SessionFromContext(ctx)
	if ok {
		response.BadRequest(rw, r, "Already logged")
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
	if len(data.Password) == 0 {
		valid.AddErrors("Password must not be empty.")
	}
	if data.Password != data.Password2 {
		valid.AddErrors("Passwords don't match.")
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
			Email:       data.Email,
			Name:        data.Name,
			IsSuperuser: false,
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

func newHandler(
	memberService MemberService,
	passwordService PasswordService,
) *handler {
	return &handler{
		memberService:   memberService,
		passwordService: passwordService,
	}
}

func Setup(
	router *mux.Router,
	memberService MemberService,
	passwordService PasswordService,
) {
	router.
		Name(RouteName).
		Methods(http.MethodPost).
		Handler(newHandler(memberService, passwordService))
}
