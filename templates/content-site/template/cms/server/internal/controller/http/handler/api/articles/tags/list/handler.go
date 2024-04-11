package list

import (
	"cms/internal/controller/http/response"
	"cms/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

const RouteName = "api.articles.tags.list"

type (
	responseData struct {
		Tags []tagData `json:"tags"`
	}
	tagData struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
)

type handler struct {
	tagService TagService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tags, err := h.tagService.All(ctx)
	if err != nil {
		logger.Load(ctx).Error(err, "failed to get tags")
		response.InternalServerError(rw, r)
		return
	}

	items := make([]tagData, len(tags))
	for i, tag := range tags {
		items[i] = tagData{
			ID:   string(tag.ID),
			Name: tag.Name,
			Slug: string(tag.Slug),
		}
	}

	response.JSON(rw, r, http.StatusOK, responseData{
		Tags: items,
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
		Methods(http.MethodGet).
		Handler(newHandler(tagService))
}
