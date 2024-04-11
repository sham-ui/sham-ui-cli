package slugify

import (
	"cms/internal/model"
	"cms/pkg/tracing"
	"context"

	"github.com/gosimple/slug"
)

const scopeName = "service.slugify"

type service struct{}

func (*service) SlugifyCategory(ctx context.Context, name string) model.CategorySlug {
	const op = "SlugifyCategory"

	_, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	return model.CategorySlug(slug.Make(name))
}

func (*service) SlugifyTag(ctx context.Context, name string) model.TagSlug {
	const op = "SlugifyTag"

	_, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	return model.TagSlug(slug.Make(name))
}

func (*service) SlugifyArticle(ctx context.Context, name string) model.ArticleSlug {
	const op = "SlugifyArticle"

	_, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	return model.ArticleSlug(slug.Make(name))
}

func New() *service {
	return &service{}
}
