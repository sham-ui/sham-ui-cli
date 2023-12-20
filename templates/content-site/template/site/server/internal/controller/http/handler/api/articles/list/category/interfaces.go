package category

import (
	"context"

	"site/internal/model"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name ArticlesService --inpackage --testonly
type ArticlesService interface {
	//nolint:lll
	ArticleListForCategory(ctx context.Context, categorySlug string, offset, limit int64) (*model.PaginatedArticleForCategory, error)
}
