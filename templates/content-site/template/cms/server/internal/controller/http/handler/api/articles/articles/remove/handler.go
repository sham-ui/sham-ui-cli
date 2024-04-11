package remove

import (
	"cms/internal/controller/http/response"
	"cms/internal/model"
	"cms/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

const RouteName = "api.articles.articles.remove"

const idKey = "id"

type responseData struct {
	Status string `json:"Status"`
}

type handler struct {
	articleService ArticleService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := mux.Vars(r)[idKey]
	if !ok {
		response.BadRequest(rw, r, "Empty id")
		return
	}

	if err := h.articleService.Delete(ctx, model.ArticleID(id)); err != nil {
		logger.Load(ctx).Error(err, "fail delete article")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusOK, &responseData{
		Status: "Article deleted",
	})
}

func newHandler(articleService ArticleService) *handler {
	return &handler{
		articleService: articleService,
	}
}

func Setup(router *mux.Router, articleService ArticleService) {
	router.
		Name(RouteName).
		Methods(http.MethodDelete).
		Path("/{" + idKey + "}").
		Handler(newHandler(articleService))
}
