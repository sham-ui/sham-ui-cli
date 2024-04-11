package list

import (
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name ArticleService --inpackage --testonly
type ArticleService interface {
	FindShortInfo(ctx context.Context, offset, limit int64) ([]model.ArticleShortInfo, error)
	Total(ctx context.Context) (int, error)
}
