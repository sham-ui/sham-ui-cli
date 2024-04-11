package password

import (
	"net/http"
	"{{ shortName }}/internal/controller/http/request"
	"{{ shortName }}/internal/controller/http/response"
	"{{ shortName }}/pkg/logger"
	"{{ shortName }}/pkg/validation"

	"github.com/gorilla/mux"
)

const RouteName = "api.members.password.update"

type requestData struct {
	NewPassword1 string `json:"NewPassword1"`
	NewPassword2 string `json:"NewPassword2"`
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

	sess, ok := request.SessionFromContext(ctx)
	if !ok {
		response.WithError(rw, r, http.StatusUnauthorized, "not authenticated")
		return
	}

	data, err := request.DecodeJSON[requestData](r)
	if err != nil {
		response.BadRequest(rw, r, "Invalid JSON")
		return
	}

	valid := validation.New()
	if len(data.NewPassword1) == 0 || len(data.NewPassword2) == 0 {
		valid.AddErrors("Password must have more than 0 characters.")
	}
	if data.NewPassword1 != data.NewPassword2 {
		valid.AddErrors("Passwords don't match.")
	}
	if valid.HasErrors() {
		response.BadRequest(rw, r, valid.Errors...)
		return
	}

	hashed, err := h.passwordService.Hash(ctx, data.NewPassword1)
	if err != nil {
		logger.Load(ctx).Error(err, "fail hash password")
		response.InternalServerError(rw, r)
		return
	}

	if err = h.memberService.UpdatePassword(ctx, sess.MemberID, hashed); err != nil {
		logger.Load(ctx).Error(err, "fail update member password")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusOK, responseData{Status: "Password updated"})
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
		Methods(http.MethodPut).
		Path("/password").
		Handler(newHandler(memberService, passwordService))
}
