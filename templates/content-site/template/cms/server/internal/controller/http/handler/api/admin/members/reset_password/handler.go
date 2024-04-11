package reset_password

import (
	"cms/internal/controller/http/request"
	"cms/internal/controller/http/response"
	"cms/internal/model"
	"cms/pkg/logger"
	"cms/pkg/validation"
	"net/http"

	"github.com/gorilla/mux"
)

const RouteName = "api.admin.members.reset_password"

const idKey = "id"

type requestData struct {
	Password1 string `json:"pass1"`
	Password2 string `json:"pass2"`
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

	valid := validation.New()
	if len(data.Password1) == 0 || len(data.Password2) == 0 {
		valid.AddErrors("Password must have more than 0 characters.")
	}
	if data.Password1 != data.Password2 {
		valid.AddErrors("Passwords don't match.")
	}
	if valid.HasErrors() {
		response.BadRequest(rw, r, valid.Errors...)
		return
	}

	hashed, err := h.passwordService.Hash(ctx, data.Password1)
	if err != nil {
		logger.Load(ctx).Error(err, "fail hash password")
		response.InternalServerError(rw, r)
		return
	}

	if err = h.memberService.UpdatePassword(ctx, model.MemberID(id), hashed); err != nil {
		logger.Load(ctx).Error(err, "fail update member password")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusOK, responseData{Status: "Password updated"})
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
		Methods(http.MethodPut).
		Path("/{" + idKey + ":[0-9]+}/password").
		Handler(newHandler(memberService, passwordService))
}
