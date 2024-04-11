package email

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

const RouteName = "api.members.email.update"

type requestData struct {
	NewEmail1 string `json:"NewEmail1"`
	NewEmail2 string `json:"NewEmail2"`
}

type responseData struct {
	Status string `json:"Status"`
}

type handler struct {
	sessionService SessionService
	memberService  MemberService
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
	data.NewEmail1 = strings.TrimSpace(data.NewEmail1)
	data.NewEmail2 = strings.TrimSpace(data.NewEmail2)

	valid := validation.New()
	if len(data.NewEmail1) == 0 || len(data.NewEmail2) == 0 {
		valid.AddErrors("Email must have more than 0 characters.")
	}
	if data.NewEmail1 != data.NewEmail2 {
		valid.AddErrors("Emails don't match.")
	}
	if valid.HasErrors() {
		response.BadRequest(rw, r, valid.Errors...)
		return
	}

	err = h.memberService.UpdateEmail(ctx, sess.MemberID, data.NewEmail1)
	switch {
	case errors.Is(err, model.ErrMemberEmailAlreadyExists):
		response.BadRequest(rw, r, "Email is already in use.")
		return
	case err != nil:
		logger.Load(ctx).Error(err, "fail update member email")
		response.InternalServerError(rw, r)
		return
	}

	if err := h.sessionService.UpdateEmail(rw, r, data.NewEmail1); err != nil {
		logger.Load(ctx).Error(err, "fail update session email")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusOK, responseData{Status: "Email updated"})
}

func newHandler(sessionService SessionService, memberService MemberService) *handler {
	return &handler{
		sessionService: sessionService,
		memberService:  memberService,
	}
}

func Setup(router *mux.Router, sessionService SessionService, memberService MemberService) {
	router.
		Name(RouteName).
		Methods(http.MethodPut).
		Path("/email").
		Handler(newHandler(sessionService, memberService))
}
