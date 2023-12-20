package detail

import (
	"errors"
	"net/http"

	"site/internal/controller/http/response"
	"site/internal/model"
	"site/pkg/logger"

	"github.com/gorilla/mux"
)

const slugKey = "slug"

type (
	responseData struct {
		Title        string   `json:"title"`
		Slug         string   `json:"slug"`
		Category     category `json:"category"`
		Tags         []tag    `json:"tags"`
		ShortContent string   `json:"shortContent"`
		Content      string   `json:"content"`
		CreatedAt    string   `json:"createdAt"`
	}

	category struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}

	tag struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
)

type handler struct {
	service ArticlesService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log := logger.Load(r.Context())
	slug := mux.Vars(r)[slugKey]
	log = log.WithValues("slug", slug)
	article, err := h.service.Article(r.Context(), slug)
	if errors.Is(err, model.ArticleNotFoundError{}) {
		response.NotFound(rw, r, "article not found")
		return
	}
	if err != nil {
		response.InternalServerError(rw, r)
		log.Error(err, "can't get article")
		return
	}
	tags := make([]tag, len(article.Tags))
	for i, item := range article.Tags {
		tags[i] = tag{
			Name: item.Name,
			Slug: item.Slug,
		}
	}
	response.JSON(rw, r, http.StatusOK, responseData{
		Title: article.Title,
		Slug:  article.Slug,
		Category: category{
			Name: article.Category.Name,
			Slug: article.Category.Slug,
		},
		Tags:         tags,
		ShortContent: article.ShortContent,
		Content:      article.Content,
		CreatedAt:    article.PublishedAt.Format("2006-01-02 15:04:05.999999999 -0700 MST"),
	})
}

func newHandler(service ArticlesService) *handler {
	return &handler{
		service: service,
	}
}

func Setup(router *mux.Router, service ArticlesService) {
	router.
		Methods("GET").
		Path("/articles/{" + slugKey + "}").
		Handler(newHandler(service))
}
