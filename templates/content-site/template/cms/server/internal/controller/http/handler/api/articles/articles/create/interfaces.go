package create

import (
	"cms/internal/controller/http/handler/api/articles/articles"
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name ArticleService --inpackage --testonly
type ArticleService interface {
	Create(ctx context.Context, data model.Article, tags []model.Tag) (model.ArticleID, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name SlugifyService --inpackage --testonly
type SlugifyService interface {
	articles.SlugifyService
}
