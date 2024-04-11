package detail

import (
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name ArticleService --inpackage --testonly
type ArticleService interface {
	FindBySlug(ctx context.Context, slug model.ArticleSlug) (model.Article, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name CategoryService --inpackage --testonly
type CategoryService interface {
	GetByID(ctx context.Context, id model.CategoryID) (*model.Category, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name ArticleTagService --inpackage --testonly
type ArticleTagService interface {
	GetTagForArticle(ctx context.Context, articleID model.ArticleID) ([]model.Tag, error)
}
