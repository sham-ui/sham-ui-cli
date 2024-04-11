package update

import (
	"cms/internal/controller/http/handler/api/articles/categories"
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name CategoryService --inpackage --testonly
type CategoryService interface {
	Update(ctx context.Context, data model.Category) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name SlugifyService --inpackage --testonly
type SlugifyService interface {
	categories.SlugifyService
}
