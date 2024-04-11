package tag

import (
	"net/http"
	"strings"

	"site/internal/controller/http/handler/api/articles/list"
	"site/internal/controller/http/request"
	"site/internal/controller/http/response"
	"site/pkg/logger"

	"github.com/gorilla/mux"
)

const RouteName = "api.articles.list.tag"

const (
	tagKey = "tag"
)

type responseData struct {
	Articles []list.Article `json:"articles"`
	Meta     meta           `json:"meta"`
}

type meta struct {
	list.Meta
	Tag string `json:"tag"`
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
	tag := strings.TrimSpace(mux.Vars(r)[tagKey])
	if tag == "" {
		response.BadRequest(rw, r, "empty tag")
		return
	}
	paginatedArticles, err := h.service.ArticleListForTag(
		r.Context(),
		tag,
		paginationParams.Offset,
		paginationParams.Limit,
	)
	if err != nil {
		response.InternalServerError(rw, r)
		log.Error(err, "can't get articles")
		return
	}
	response.JSON(rw, r, http.StatusOK, responseData{
		Articles: list.ArticlesFromModel(paginatedArticles.Articles),
		Meta: meta{
			Meta: list.Meta{
				Offset: paginationParams.Offset,
				Limit:  paginationParams.Limit,
				Total:  paginatedArticles.Total,
			},
			Tag: paginatedArticles.TagName,
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
		Queries("tag", "{"+tagKey+"}").
		Path("/articles").
		Handler(newHandler(service))
}
