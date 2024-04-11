package list

import (
	"cms/internal/controller/http/request"
	"cms/internal/controller/http/response"
	"cms/pkg/logger"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const RouteName = "api.articles.articles.list"

type (
	responseData struct {
		Articles []article `json:"articles"`
		Meta     meta      `json:"meta"`
	}
	article struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		Slug        string    `json:"slug"`
		CategoryID  string    `json:"category_id"`
		PublishedAt time.Time `json:"published_at"`
	}
	meta struct {
		Offset int64 `json:"offset"`
		Limit  int64 `json:"limit"`
		Total  int64 `json:"total"`
	}
)

type handler struct {
	articleService ArticleService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pagination, ok := request.ExtractPaginationParams(rw, r)
	if !ok {
		return
	}

	articles, err := h.articleService.FindShortInfo(ctx, pagination.Offset, pagination.Limit)
	if err != nil {
		logger.Load(ctx).Error(err, "failed to find articles")
		response.InternalServerError(rw, r)
		return
	}

	total, err := h.articleService.Total(ctx)
	if err != nil {
		logger.Load(ctx).Error(err, "failed to count articles")
		response.InternalServerError(rw, r)
		return
	}

	articlesResponse := make([]article, len(articles))
	for i, art := range articles {
		articlesResponse[i] = article{
			ID:          string(art.ID),
			Title:       art.Title,
			Slug:        string(art.Slug),
			CategoryID:  string(art.CategoryID),
			PublishedAt: art.PublishedAt,
		}
	}

	response.JSON(rw, r, http.StatusOK, responseData{
		Articles: articlesResponse,
		Meta: meta{
			Offset: pagination.Offset,
			Limit:  pagination.Limit,
			Total:  int64(total),
		},
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
		Methods(http.MethodGet).
		Handler(newHandler(articleService))
}
