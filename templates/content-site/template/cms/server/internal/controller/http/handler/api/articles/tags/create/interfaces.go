package create

import (
	"cms/internal/controller/http/handler/api/articles/tags"
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name TagService --inpackage --testonly
type TagService interface {
	Create(ctx context.Context, data model.Tag) (model.TagID, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name SlugifyService --inpackage --testonly
type SlugifyService interface {
	tags.SlugifyService
}
