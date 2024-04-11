package tag

import (
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name ArticleService --inpackage --testonly
type ArticleService interface {
	FindShortInfoWithCategoryForTag(
		ctx context.Context,
		categorySlug model.TagSlug,
		offset, limit int64,
	) ([]model.ArticleShortInfoWithCategory, error)
	TotalForTag(ctx context.Context, slug model.TagSlug) (int, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name TagService --inpackage --testonly
type TagService interface {
	GetBySlug(ctx context.Context, slug model.TagSlug) (model.Tag, error)
}
