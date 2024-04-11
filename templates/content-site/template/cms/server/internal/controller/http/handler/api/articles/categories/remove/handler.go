package remove

import (
	"cms/internal/controller/http/response"
	"cms/internal/model"
	"cms/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

const RouteName = "api.articles.categories.remove"

const idKey = "id"

type responseData struct {
	Status string `json:"Status"`
}

type handler struct {
	categoryService CategoryService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := mux.Vars(r)[idKey]
	if !ok {
		response.BadRequest(rw, r, "Empty id")
		return
	}

	if err := h.categoryService.Delete(ctx, model.CategoryID(id)); err != nil {
		logger.Load(ctx).Error(err, "fail delete category")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusOK, &responseData{
		Status: "Category deleted",
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
		Methods(http.MethodDelete).
		Path("/{" + idKey + "}").
		Handler(newHandler(categoryService))
}
