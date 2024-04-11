package articles

import (
	"cms/internal/controller/http/request"
	"cms/internal/controller/http/response"
	"cms/internal/model"
	"cms/pkg/validation"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name SlugifyService --inpackage --testonly --with-expecter
type SlugifyService interface {
	SlugifyArticle(ctx context.Context, name string) model.ArticleSlug
	SlugifyTag(ctx context.Context, name string) model.TagSlug
}

type (
	requestData struct {
		Title       string       `json:"title"`
		CategoryID  string       `json:"category_id"`
		ShortBody   string       `json:"short_body"`
		Body        string       `json:"body"`
		PublishedAt time.Time    `json:"published_at"`
		Tags        []articleTag `json:"tags"`
	}
	articleTag struct {
		Slug string `json:"slug"`
		Name string `json:"name"`
	}
)

func ExtractAndValidateData(
	slugger SlugifyService,
	rw http.ResponseWriter,
	r *http.Request,
) (*model.Article, []model.Tag, bool) {
	ctx := r.Context()

	data, err := request.DecodeJSON[requestData](r)
	if err != nil {
		response.BadRequest(rw, r, "Invalid JSON")
		return nil, nil, false
	}
	data.Title = strings.TrimSpace(data.Title)
	data.ShortBody = strings.TrimSpace(data.ShortBody)
	data.Body = strings.TrimSpace(data.Body)

	valid := validation.New()
	if len(data.Title) == 0 {
		valid.AddErrors("Title must not be empty.")
	}
	if len(data.CategoryID) == 0 {
		valid.AddErrors("Category must not be empty.")
	}
	if valid.HasErrors() {
		response.BadRequest(rw, r, valid.Errors...)
		return nil, nil, false
	}

	slug := slugger.SlugifyArticle(ctx, data.Title)

	tags := make([]model.Tag, len(data.Tags))
	for i, tag := range data.Tags {
		item := model.Tag{ //nolint:exhaustruct
			Name: tag.Name,
		}
		if tag.Slug == "" {
			item.Name = strings.TrimSpace(tag.Name)
			item.Slug = slugger.SlugifyTag(ctx, item.Name)
		} else {
			item.Slug = model.TagSlug(tag.Slug)
		}
		tags[i] = item
	}

	return &model.Article{ //nolint:exhaustruct
		Slug:        slug,
		Title:       data.Title,
		CategoryID:  model.CategoryID(data.CategoryID),
		ShortBody:   data.ShortBody,
		Body:        data.Body,
		PublishedAt: data.PublishedAt,
	}, tags, true
}

func HandleError(err error, rw http.ResponseWriter, r *http.Request) bool {
	switch {
	case errors.Is(err, model.ErrArticleTitleAlreadyExists):
		response.BadRequest(rw, r, "Title is already in use.")
		return true
	case errors.Is(err, model.ErrArticleSlugAlreadyExists):
		response.BadRequest(rw, r, "Slug is already in use.")
		return true
	case errors.Is(err, model.ErrCategoryNotFound):
		response.BadRequest(rw, r, "Category not found.")
		return true
	}
	return false
}
