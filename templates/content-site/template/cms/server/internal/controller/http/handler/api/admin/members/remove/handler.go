package remove

import (
	"cms/internal/controller/http/response"
	"cms/internal/model"
	"cms/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

const RouteName = "api.admin.members.remove"

const idKey = "id"

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

	if err := h.memberService.Delete(ctx, model.MemberID(id)); err != nil {
		logger.Load(ctx).Error(err, "fail delete member")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusOK, &responseData{
		Status: "Member deleted",
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
		Methods(http.MethodDelete).
		Path("/{" + idKey + ":[0-9]+}").
		Handler(newHandler(memberService))
}
