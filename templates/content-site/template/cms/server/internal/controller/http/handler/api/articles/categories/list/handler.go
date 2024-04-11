package list

import (
	"cms/internal/controller/http/response"
	"cms/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

const RouteName = "api.articles.categories.list"

type (
	responseData struct {
		Categories []categoryData `json:"categories"`
	}
	categoryData struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
)

type handler struct {
	categoryService CategoryService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := h.categoryService.All(ctx)
	if err != nil {
		logger.Load(ctx).Error(err, "failed to get categories")
		response.InternalServerError(rw, r)
		return
	}

	items := make([]categoryData, len(categories))
	for i, category := range categories {
		items[i] = categoryData{
			ID:   string(category.ID),
			Name: category.Name,
			Slug: string(category.Slug),
		}
	}

	response.JSON(rw, r, http.StatusOK, responseData{
		Categories: items,
	})
}

func newHandler(categoryService CategoryService) *handler {
	return &handler{
		categoryService: categoryService,
	}
}

func Setup(router *mux.Router, categoryService CategoryService) {
	router.
		Name(RouteName).
		Methods(http.MethodGet).
		Handler(newHandler(categoryService))
}
