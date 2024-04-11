package update

import (
	"cms/internal/controller/http/handler/api/articles/tags"
	"cms/internal/controller/http/response"
	"cms/internal/model"
	"cms/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

const RouteName = "api.articles.tags.update"

const idKey = "id"

type responseData struct {
	Status string `json:"Status"`
}

type handler struct {
	tagService     TagService
	slugifyService SlugifyService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := mux.Vars(r)[idKey]
	if !ok {
		response.BadRequest(rw, r, "Empty id")
		return
	}
	tag, valid := tags.ExtractAndValidateData(h.slugifyService, rw, r)
	if !valid {
		return
	}
	tag.ID = model.TagID(id)

	err := h.tagService.Update(ctx, *tag)
	handled := tags.HandleError(err, rw, r)
	switch {
	case handled:
		return
	case err != nil:
		logger.Load(ctx).Error(err, "failed to update tag")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusOK, &responseData{
		Status: "Tag updated",
	})
}

func newHandler(tagService TagService, slugifyService SlugifyService) *handler {
	return &handler{
		tagService:     tagService,
		slugifyService: slugifyService,
	}
}

func Setup(router *mux.Router, tagService TagService, slugifyService SlugifyService) {
	router.
		Name(RouteName).
		Methods(http.MethodPut).
		Path("/{" + idKey + ":[0-9]+}").
		Handler(newHandler(tagService, slugifyService))
}
