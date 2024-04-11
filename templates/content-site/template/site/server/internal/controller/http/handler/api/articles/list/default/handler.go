package _default

import (
	"net/http"

	"site/internal/controller/http/handler/api/articles/list"
	"site/internal/controller/http/request"
	"site/internal/controller/http/response"
	"site/pkg/logger"

	"github.com/gorilla/mux"
)

const RouteName = "api.articles.list.default"

type responseData struct {
	Articles []list.Article `json:"articles"`
	Meta     list.Meta      `json:"meta"`
}

type handler struct {
	service ArticlesService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log := logger.Load(r.Context())
	paginationParams, valid := request.ExtractPaginationParams(rw, r)
	if !valid {
		return
	}
	paginatedArticles, err := h.service.Articles(r.Context(), paginationParams.Offset, paginationParams.Limit)
	if err != nil {
		response.InternalServerError(rw, r)
		log.Error(err, "can't get articles")
		return
	}
	response.JSON(rw, r, http.StatusOK, responseData{
		Articles: list.ArticlesFromModel(paginatedArticles.Articles),
		Meta: list.Meta{
			Offset: paginationParams.Offset,
			Limit:  paginationParams.Limit,
			Total:  paginatedArticles.Total,
		},
	})
}

func newHandler(service ArticlesService) *handler {
	return &handler{
		service: service,
	}
}

func Setup(router *mux.Router, service ArticlesService) {
	router.
		Name(RouteName).
		Methods("GET").
		Path("/articles").
		Handler(newHandler(service))
}
