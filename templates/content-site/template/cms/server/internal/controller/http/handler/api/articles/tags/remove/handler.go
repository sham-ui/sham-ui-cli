package remove

import (
	"cms/internal/controller/http/response"
	"cms/internal/model"
	"cms/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

const RouteName = "api.articles.tags.remove"

const idKey = "id"

type responseData struct {
	Status string `json:"Status"`
}

type handler struct {
	tagService TagService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := mux.Vars(r)[idKey]
	if !ok {
		response.BadRequest(rw, r, "Empty id")
		return
	}

	if err := h.tagService.Delete(ctx, model.TagID(id)); err != nil {
		logger.Load(ctx).Error(err, "fail delete tag")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusOK, &responseData{
		Status: "Tag deleted",
	})
}

func newHandler(tagService TagService) *handler {
	return &handler{
		tagService: tagService,
	}
}

func Setup(router *mux.Router, tagService TagService) {
	router.
		Name(RouteName).
		Methods(http.MethodDelete).
		Path("/{" + idKey + "}").
		Handler(newHandler(tagService))
}
