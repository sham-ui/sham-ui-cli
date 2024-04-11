package detail

import (
	"cms/internal/controller/http/response"
	"cms/internal/model"
	"cms/pkg/logger"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const RouteName = "api.articles.articles.detail"

const idKey = "id"

type (
	responseData struct {
		Title       string    `json:"title"`
		Slug        string    `json:"slug"`
		CategoryID  string    `json:"category_id"`
		ShortBody   string    `json:"short_body"`
		Body        string    `json:"body"`
		PublishedAt time.Time `json:"published_at"`
		Tags        []string  `json:"tags"`
	}
)

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

	article, err := h.articleService.FindByID(ctx, model.ArticleID(id))
	switch {
	case errors.Is(err, model.ErrArticleNotFound):
		response.NotFound(rw, r, "Article not found")
		return
	case err != nil:
		logger.Load(ctx).Error(err, "fail find article")
		response.InternalServerError(rw, r)
		return
	}

	tags, err := h.articleService.GetTags(ctx, model.ArticleID(id))
	if err != nil {
		logger.Load(ctx).Error(err, "fail get article tags")
		response.InternalServerError(rw, r)
		return
	}
	tagSlugs := make([]string, len(tags))
	for i, t := range tags {
		tagSlugs[i] = string(t.Slug)
	}

	response.JSON(
		rw, r, http.StatusOK, &responseData{
			Title:       article.Title,
			Slug:        string(article.Slug),
			CategoryID:  string(article.CategoryID),
			ShortBody:   article.ShortBody,
			Body:        article.Body,
			PublishedAt: article.PublishedAt,
			Tags:        tagSlugs,
		},
	)
}

func newHandler(articleService ArticleService) *handler {
	return &handler{
		articleService: articleService,
	}
}

func Setup(router *mux.Router, articleService ArticleService) {
	router.
		Name(RouteName).
		Methods(http.MethodGet).
		Path("/{" + idKey + ":[0-9]+}").
		Handler(newHandler(articleService))
}
