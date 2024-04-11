package create

import (
	"cms/internal/controller/http/handler/api/articles/categories"
	"cms/internal/controller/http/response"
	"cms/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

const RouteName = "api.articles.categories.create"

type responseData struct {
	Status string `json:"Status"`
}

type handler struct {
	categoryService CategoryService
	slugifyService  SlugifyService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	category, valid := categories.ExtractAndValidateData(h.slugifyService, rw, r)
	if !valid {
		return
	}

	err := h.categoryService.Create(ctx, *category)
	handled := categories.HandleError(err, rw, r)
	switch {
	case handled:
		return
	case err != nil:
		logger.Load(ctx).Error(err, "failed to create category")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusCreated, &responseData{
		Status: "Category created",
	})
}

func newHandler(categoryService CategoryService, slugifyService SlugifyService) *handler {
	return &handler{
		categoryService: categoryService,
		slugifyService:  slugifyService,
	}
}

func Setup(router *mux.Router, categoryService CategoryService, slugifyService SlugifyService) {
	router.
		Name(RouteName).
		Methods(http.MethodPost).
		Handler(newHandler(categoryService, slugifyService))
}
