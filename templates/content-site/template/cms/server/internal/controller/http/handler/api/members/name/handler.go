package name

import (
	"cms/internal/controller/http/request"
	"cms/internal/controller/http/response"
	"cms/pkg/logger"
	"cms/pkg/validation"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const RouteName = "api.members.name.update"

type requestData struct {
	NewName string `json:"NewName"`
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
	data.NewName = strings.TrimSpace(data.NewName)

	valid := validation.New()
	if len(data.NewName) == 0 {
		valid.AddErrors("Name must have more than 0 characters.")
	}
	if valid.HasErrors() {
		response.BadRequest(rw, r, valid.Errors...)
		return
	}

	if err = h.memberService.UpdateName(ctx, sess.MemberID, data.NewName); err != nil {
		logger.Load(ctx).Error(err, "fail update member name")
		response.InternalServerError(rw, r)
		return
	}

	if err := h.sessionService.UpdateName(rw, r, data.NewName); err != nil {
		logger.Load(ctx).Error(err, "fail update session name")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusOK, responseData{Status: "Name updated"})
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
		Path("/name").
		Handler(newHandler(sessionService, memberService))
}
