package update

import (
	"cms/internal/controller/http/handler/api/articles/articles"
	"cms/internal/controller/http/response"
	"cms/internal/model"
	"cms/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

const RouteName = "api.articles.articles.update"

const idKey = "id"

type responseData struct {
	Status string `json:"Status"`
}

type handler struct {
	articleService ArticleService
	slugifyService SlugifyService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := mux.Vars(r)[idKey]
	if !ok {
		response.BadRequest(rw, r, "Empty id")
		return
	}
	article, tags, valid := articles.ExtractAndValidateData(h.slugifyService, rw, r)
	if !valid {
		return
	}
	article.ID = model.ArticleID(id)

	err := h.articleService.Update(ctx, *article, tags)
	handled := articles.HandleError(err, rw, r)
	switch {
	case handled:
		return
	case err != nil:
		logger.Load(ctx).Error(err, "failed to update article")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusOK, &responseData{
		Status: "Article updated",
	})
}

func newHandler(articleService ArticleService, slugifyService SlugifyService) *handler {
	return &handler{
		articleService: articleService,
		slugifyService: slugifyService,
	}
}

func Setup(router *mux.Router, articleService ArticleService, slugifyService SlugifyService) {
	router.
		Name(RouteName).
		Methods(http.MethodPut).
		Path("/{" + idKey + ":[0-9]+}").
		Handler(newHandler(articleService, slugifyService))
}
