package update

import (
	"cms/internal/controller/http/handler/api/articles/tags"
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name TagService --inpackage --testonly
type TagService interface {
	Update(ctx context.Context, data model.Tag) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name SlugifyService --inpackage --testonly
type SlugifyService interface {
	tags.SlugifyService
}
