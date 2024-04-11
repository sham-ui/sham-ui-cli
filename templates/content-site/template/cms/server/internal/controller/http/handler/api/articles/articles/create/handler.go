package create

import (
	"cms/internal/controller/http/handler/api/articles/articles"
	"cms/internal/controller/http/response"
	"cms/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

const RouteName = "api.articles.articles.create"

type responseData struct {
	Status string `json:"Status"`
}

type handler struct {
	articleService ArticleService
	slugifyService SlugifyService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	article, tags, valid := articles.ExtractAndValidateData(h.slugifyService, rw, r)
	if !valid {
		return
	}

	_, err := h.articleService.Create(ctx, *article, tags)
	handled := articles.HandleError(err, rw, r)
	switch {
	case handled:
		return
	case err != nil:
		logger.Load(ctx).Error(err, "failed to create article")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusCreated, &responseData{
		Status: "Article created",
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
		Methods(http.MethodPost).
		Handler(newHandler(articleService, slugifyService))
}
