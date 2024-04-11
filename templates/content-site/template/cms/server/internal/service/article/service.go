package article

import (
	"cms/internal/model"
	"cms/pkg/set"
	"cms/pkg/tracing"
	"cms/pkg/transactor"
	"context"
	"errors"
	"fmt"
)

const scopeName = "service.article"

type service struct {
	transactor           transactor.Transactor
	articleRepository    ArticleRepository
	tagRepository        TagRepository
	articleTagRepository ArticleTagRepository
}

func (s *service) FindShortInfo(ctx context.Context, offset, limit int64) ([]model.ArticleShortInfo, error) {
	const op = "FindShortInfo"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	res, err := s.articleRepository.FindShortInfo(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("find short info: %w", err)
	}
	return res, nil
}

func (s *service) FindByID(ctx context.Context, id model.ArticleID) (model.Article, error) {
	const op = "FindByID"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	res, err := s.articleRepository.FindByID(ctx, id)
	if err != nil {
		return res, fmt.Errorf("find by id: %w", err)
	}
	return res, nil
}

func (s *service) Total(ctx context.Context) (int, error) {
	const op = "Total"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	res, err := s.articleRepository.Total(ctx)
	if err != nil {
		return 0, fmt.Errorf("total: %w", err)
	}
	return res, nil
}

func (s *service) GetTags(ctx context.Context, articleID model.ArticleID) ([]model.Tag, error) {
	const op = "GetTags"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	tags, err := s.articleTagRepository.GetTagForArticle(ctx, articleID)
	if err != nil {
		return nil, fmt.Errorf("get tags for article: %w", err)
	}
	return tags, nil
}

func (s *service) Create(ctx context.Context, data model.Article, tags []model.Tag) (model.ArticleID, error) {
	const op = "Create"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	var id model.ArticleID
	if err := s.transactor.WithinTransaction(
		ctx,
		func(ctx context.Context) error {
			var err error
			id, err = s.articleRepository.Create(ctx, data)
			if err != nil {
				return fmt.Errorf("create article: %w", err)
			}

			tagsIDs, err := s.getOrCreateTags(ctx, tags)
			if err != nil {
				return fmt.Errorf("get or create tags: %w", err)
			}

			for _, tagId := range tagsIDs {
				if err := s.articleTagRepository.Create(ctx, id, tagId); err != nil {
					return fmt.Errorf("create article tag(tag_id=%s, article_id=%s): %w", tagId, id, err)
				}
			}

			return nil
		},
	); err != nil {
		return "", fmt.Errorf("within transaction: %w", err)
	}
	return id, nil
}

func (s *service) Update(ctx context.Context, data model.Article, tags []model.Tag) error {
	const op = "Update"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	if err := s.transactor.WithinTransaction(
		ctx, func(ctx context.Context) error {
			if err := s.articleRepository.Update(ctx, data); err != nil {
				return fmt.Errorf("update article: %w", err)
			}

			existedTagIDs, err := s.articleTagRepository.GetTagIDs(ctx, data.ID)
			if err != nil {
				return fmt.Errorf("get article tag ids(article_id=%s): %w", data.ID, err)
			}

			tagsIDs, err := s.getOrCreateTags(ctx, tags)
			if err != nil {
				return fmt.Errorf("get or create tags: %w", err)
			}

			addIDs, deleteIDs := set.New(existedTagIDs).Difference(set.New(tagsIDs))
			for _, tagId := range addIDs {
				if err := s.articleTagRepository.Create(ctx, data.ID, tagId); err != nil {
					return fmt.Errorf("create article tag(tag_id=%s, article_id=%s): %w", tagId, data.ID, err)
				}
			}
			for _, tagId := range deleteIDs {
				if err := s.articleTagRepository.Delete(ctx, data.ID, tagId); err != nil {
					return fmt.Errorf("delete article tag(tag_id=%s, article_id=%s): %w", tagId, data.ID, err)
				}
			}

			return nil
		},
	); err != nil {
		return fmt.Errorf("within transaction: %w", err)
	}
	return nil
}

func (s *service) Delete(ctx context.Context, id model.ArticleID) error {
	const op = "Delete"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	if err := s.articleRepository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete article: %w", err)
	}
	return nil
}

func (s *service) getOrCreateTags(ctx context.Context, tags []model.Tag) ([]model.TagID, error) {
	const op = "getOrCreateTags"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	tagsIDs := make([]model.TagID, len(tags))
	for i, tag := range tags {
		rec, err := s.tagRepository.GetBySlug(ctx, tag.Slug)
		switch {
		case errors.Is(err, model.ErrTagNotFound):
			id, err := s.tagRepository.Create(ctx, tag)
			if err != nil {
				return nil, fmt.Errorf("create tag(slug=%s): %w", tag.Slug, err)
			}
			tagsIDs[i] = id
		case err != nil:
			return nil, fmt.Errorf("get tag by slug(slug=%s): %w", tag.Slug, err)
		default:
			tagsIDs[i] = rec.ID
		}
	}
	return tagsIDs, nil
}

func New(
	trans transactor.Transactor,
	articleRepo ArticleRepository,
	tagRepo TagRepository,
	articleTagRepo ArticleTagRepository,
) *service {
	return &service{
		transactor:           trans,
		articleRepository:    articleRepo,
		tagRepository:        tagRepo,
		articleTagRepository: articleTagRepo,
	}
}
