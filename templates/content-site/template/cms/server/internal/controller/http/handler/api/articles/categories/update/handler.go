package update

import (
	"cms/internal/controller/http/handler/api/articles/categories"
	"cms/internal/controller/http/response"
	"cms/internal/model"
	"cms/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

const RouteName = "api.articles.categories.update"

const idKey = "id"

type responseData struct {
	Status string `json:"Status"`
}

type handler struct {
	categoryService CategoryService
	slugifyService  SlugifyService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := mux.Vars(r)[idKey]
	if !ok {
		response.BadRequest(rw, r, "Empty id")
		return
	}
	category, valid := categories.ExtractAndValidateData(h.slugifyService, rw, r)
	if !valid {
		return
	}
	category.ID = model.CategoryID(id)

	err := h.categoryService.Update(ctx, *category)
	handled := categories.HandleError(err, rw, r)
	switch {
	case handled:
		return
	case err != nil:
		logger.Load(ctx).Error(err, "failed to update category")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusOK, &responseData{
		Status: "Category updated",
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
		Methods(http.MethodPut).
		Path("/{" + idKey + ":[0-9]+}").
		Handler(newHandler(categoryService, slugifyService))
}
