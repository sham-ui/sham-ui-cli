package login

import (
	"errors"
	"net/http"
	"{{ shortName }}/internal/controller/http/request"
	"{{ shortName }}/internal/controller/http/response"
	"{{ shortName }}/internal/model"
	"{{ shortName }}/pkg/logger"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

const RouteName = "api.session.login"

type requestData struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type responseData struct {
	Status      string `json:"Status"`
	Name        string `json:"Name"`
	Email       string `json:"Email"`
	IsSuperuser bool   `json:"IsSuperuser"`
}

type handler struct {
	sessionService    SessionService
	memberService     MemberService
	passwordService   PasswordService
	csrfRequestHeader string
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, ok := request.SessionFromContext(ctx)
	if ok {
		response.BadRequest(rw, r, "Already logged in")
		return
	}

	data, err := request.DecodeJSON[requestData](r)
	if err != nil {
		response.BadRequest(rw, r, "Invalid JSON")
		return
	}

	memberWithPassword, err := h.memberService.GetByEmail(ctx, data.Email)
	switch {
	case errors.Is(err, model.ErrMemberNotFound):
		response.BadRequest(rw, r, "Member not found")
		return
	case err != nil:
		logger.Load(ctx).Error(err, "fail get member")
		response.InternalServerError(rw, r)
		return
	}

	if err := h.passwordService.Compare(ctx, memberWithPassword.HashedPassword, data.Password); err != nil {
		response.BadRequest(rw, r, "Incorrect username or password")
		return
	}

	if err := h.sessionService.Create(rw, r, &memberWithPassword.Member); err != nil {
		logger.Load(ctx).Error(err, "fail create session")
		response.InternalServerError(rw, r)
		return
	}

	rw.Header().Set(h.csrfRequestHeader, csrf.Token(r))

	response.JSON(rw, r, http.StatusOK, &responseData{
		Status:      "OK",
		Name:        memberWithPassword.Name,
		Email:       memberWithPassword.Email,
		IsSuperuser: memberWithPassword.IsSuperuser,
	})
}

func newHandler(
	sessionService SessionService,
	memberService MemberService,
	passwordService PasswordService,
	csrfRequestHeader string,
) *handler {
	return &handler{
		sessionService:    sessionService,
		memberService:     memberService,
		passwordService:   passwordService,
		csrfRequestHeader: csrfRequestHeader,
	}
}

func Setup(
	router *mux.Router,
	sessionService SessionService,
	memberService MemberService,
	passwordService PasswordService,
	csrfRequestHeader string,
) {
	router.
		Name(RouteName).
		Methods(http.MethodPost).
		Path("/login").
		Handler(newHandler(sessionService, memberService, passwordService, csrfRequestHeader))
}
