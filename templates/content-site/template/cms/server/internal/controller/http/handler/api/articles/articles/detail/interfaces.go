package detail

import (
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name ArticleService --inpackage --testonly --with-expecter
type ArticleService interface {
	FindByID(ctx context.Context, id model.ArticleID) (model.Article, error)
	GetTags(ctx context.Context, articleID model.ArticleID) ([]model.Tag, error)
}
