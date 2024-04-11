package list

import (
	"net/http"
	"{{ shortName }}/internal/controller/http/request"
	"{{ shortName }}/internal/controller/http/response"
	"{{ shortName }}/pkg/logger"

	"github.com/gorilla/mux"
)

const RouteName = "api.admin.members.list"

type (
	responseData struct {
		Members []member `json:"members"`
		Meta    meta     `json:"meta"`
	}
	member struct {
		ID          string `json:"ID"`
		Name        string `json:"Name"`
		Email       string `json:"Email"`
		IsSuperuser bool   `json:"IsSuperuser"`
	}
	meta struct {
		Offset int64 `json:"offset"`
		Limit  int64 `json:"limit"`
		Total  int64 `json:"total"`
	}
)

type handler struct {
	memberService MemberService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pagination, ok := request.ExtractPaginationParams(rw, r)
	if !ok {
		return
	}

	members, err := h.memberService.Find(ctx, pagination.Offset, pagination.Limit)
	if err != nil {
		logger.Load(ctx).Error(err, "failed to find members")
		response.InternalServerError(rw, r)
		return
	}

	total, err := h.memberService.Total(ctx)
	if err != nil {
		logger.Load(ctx).Error(err, "failed to count members")
		response.InternalServerError(rw, r)
		return
	}

	membersResponse := make([]member, len(members))
	for i, m := range members {
		membersResponse[i] = member{
			ID:          string(m.ID),
			Name:        m.Name,
			Email:       m.Email,
			IsSuperuser: m.IsSuperuser,
		}
	}

	response.JSON(rw, r, http.StatusOK, responseData{
		Members: membersResponse,
		Meta: meta{
			Offset: pagination.Offset,
			Limit:  pagination.Limit,
			Total:  int64(total),
		},
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
		Methods(http.MethodGet).
		Handler(newHandler(memberService))
}
