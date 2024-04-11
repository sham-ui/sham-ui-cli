package query

import (
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name ArticleService --inpackage --testonly
type ArticleService interface {
	FindShortInfoWithCategoryForQuery(
		ctx context.Context,
		query string,
		offset, limit int64,
	) ([]model.ArticleShortInfoWithCategory, error)
	TotalForQuery(ctx context.Context, query string) (int, error)
}
