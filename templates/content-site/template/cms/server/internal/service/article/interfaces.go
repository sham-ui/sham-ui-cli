package article

import (
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name ArticleRepository --inpackage --testonly
type ArticleRepository interface {
	Create(ctx context.Context, data model.Article) (model.ArticleID, error)
	Update(ctx context.Context, data model.Article) error
	FindShortInfo(ctx context.Context, offset, limit int64) ([]model.ArticleShortInfo, error)
	FindByID(ctx context.Context, id model.ArticleID) (model.Article, error)
	Total(ctx context.Context) (int, error)
	Delete(ctx context.Context, id model.ArticleID) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name TagRepository --inpackage --testonly
type TagRepository interface {
	GetBySlug(ctx context.Context, slug model.TagSlug) (model.Tag, error)
	Create(ctx context.Context, data model.Tag) (model.TagID, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name ArticleTagRepository --inpackage --testonly
type ArticleTagRepository interface {
	Create(ctx context.Context, articleID model.ArticleID, tagID model.TagID) error
	Delete(ctx context.Context, articleID model.ArticleID, tagID model.TagID) error
	GetTagIDs(ctx context.Context, articleID model.ArticleID) ([]model.TagID, error)
	GetTagForArticle(ctx context.Context, articleID model.ArticleID) ([]model.Tag, error)
}
