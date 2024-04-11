package create

import (
	"cms/internal/controller/http/handler/api/articles/tags"
	"cms/internal/controller/http/response"
	"cms/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

const RouteName = "api.articles.tags.create"

type responseData struct {
	Status string `json:"Status"`
}

type handler struct {
	tagService     TagService
	slugifyService SlugifyService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tag, valid := tags.ExtractAndValidateData(h.slugifyService, rw, r)
	if !valid {
		return
	}

	_, err := h.tagService.Create(ctx, *tag)
	handled := tags.HandleError(err, rw, r)
	switch {
	case handled:
		return
	case err != nil:
		logger.Load(ctx).Error(err, "failed to create tag")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusCreated, &responseData{
		Status: "Tag created",
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
		Methods(http.MethodPost).
		Handler(newHandler(tagService, slugifyService))
}
