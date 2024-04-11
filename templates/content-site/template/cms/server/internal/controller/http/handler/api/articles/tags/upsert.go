package tags

import (
	"cms/internal/controller/http/request"
	"cms/internal/controller/http/response"
	"cms/internal/model"
	"context"
	"errors"
	"net/http"
	"strings"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name SlugifyService --inpackage --testonly --with-expecter
type SlugifyService interface {
	SlugifyTag(ctx context.Context, name string) model.TagSlug
}

type upsertRequestData struct {
	Name string `json:"name"`
}

func ExtractAndValidateData(slugger SlugifyService, rw http.ResponseWriter, r *http.Request) (*model.Tag, bool) {
	data, err := request.DecodeJSON[upsertRequestData](r)
	if err != nil {
		response.BadRequest(rw, r, "Invalid JSON")
		return nil, false
	}
	data.Name = strings.TrimSpace(data.Name)

	if len(data.Name) == 0 {
		response.BadRequest(rw, r, "Name must not be empty.")
		return nil, false
	}

	slug := slugger.SlugifyTag(r.Context(), data.Name)

	return &model.Tag{ //nolint:exhaustruct
		Slug: slug,
		Name: data.Name,
	}, true
}

func HandleError(err error, rw http.ResponseWriter, r *http.Request) bool {
	switch {
	case errors.Is(err, model.ErrTagSlugAlreadyExists):
		response.BadRequest(rw, r, "Slug is already in use.")
		return true
	case errors.Is(err, model.ErrTagNameAlreadyExists):
		response.BadRequest(rw, r, "Name is already in use.")
		return true
	}
	return false
}
