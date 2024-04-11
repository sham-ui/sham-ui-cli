package category

import (
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name ArticleService --inpackage --testonly
type ArticleService interface {
	FindShortInfoWithCategoryForCategory(
		ctx context.Context,
		categorySlug model.CategorySlug,
		offset, limit int64,
	) ([]model.ArticleShortInfoWithCategory, error)
	TotalForCategory(ctx context.Context, slug model.CategorySlug) (int, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name CategoryService --inpackage --testonly
type CategoryService interface {
	GetBySlug(ctx context.Context, slug model.CategorySlug) (*model.Category, error)
}
