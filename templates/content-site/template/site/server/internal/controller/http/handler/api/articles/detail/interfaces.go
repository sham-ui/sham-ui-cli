package detail

import (
	"context"

	"site/internal/model"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name ArticlesService --inpackage --testonly
type ArticlesService interface {
	Article(ctx context.Context, slug string) (*model.Article, error)
}
