package article_tag

import (
	"cms/internal/model"
	"cms/pkg/postgres"
	"cms/pkg/tracing"
	"context"
	"fmt"
)

const scopeName = "repository.postgres.article_tag"

const (
	articleTagTagIdArticleIdUniqueConstraint = "article_tag_tag_id_article_id_unique"
	articleTagTagIDFKeyConstraint            = "article_tag_tag_id_fkey"
	articleTagArticleIDFKeyConstraint        = "article_tag_article_id_fkey"
)

type repository struct {
	db *postgres.Database
}

func (r *repository) Create(ctx context.Context, articleID model.ArticleID, tagID model.TagID) error {
	const op = "Create"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO article_tag(article_id, tag_id) VALUES ($1, $2)",
		articleID,
		tagID,
	)
	switch {
	case postgres.IsUniqueViolationError(err, articleTagTagIdArticleIdUniqueConstraint):
		return model.ErrArticleTagAlreadyExists
	case postgres.IsForeignKeyViolationError(err, articleTagTagIDFKeyConstraint):
		return model.ErrTagNotFound
	case postgres.IsForeignKeyViolationError(err, articleTagArticleIDFKeyConstraint):
		return model.ErrArticleNotFound
	case err != nil:
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func (r *repository) GetTagIDs(ctx context.Context, articleID model.ArticleID) ([]model.TagID, error) {
	const op = "GetTagIDs"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT tag_id FROM article_tag WHERE article_id = $1 ORDER BY tag_id",
		articleID,
	)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var tagIDs []model.TagID
	for rows.Next() {
		var tagID string
		if err := rows.Scan(&tagID); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		tagIDs = append(tagIDs, model.TagID(tagID))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return tagIDs, nil
}

func (r *repository) GetTagForArticle(ctx context.Context, articleID model.ArticleID) ([]model.Tag, error) {
	const op = "GetTagForArticle"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT t.id, t.name, t.slug
		FROM article_tag a_t
		JOIN tag t ON t.id = a_t.tag_id
		WHERE a_t.article_id = $1
		ORDER BY t.name`,
		articleID,
	)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		var data tag
		if err := rows.Scan(&data.ID, &data.Name, &data.Slug); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		tags = append(tags, model.Tag{
			ID:   model.TagID(data.ID),
			Slug: model.TagSlug(data.Slug),
			Name: data.Name,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return tags, nil
}

func (r *repository) Delete(ctx context.Context, articleID model.ArticleID, tagID model.TagID) error {
	const op = "Delete"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	if _, err := r.db.ExecContext(
		ctx,
		"DELETE FROM article_tag WHERE article_id = $1 AND tag_id = $2",
		articleID,
		tagID,
	); err != nil {
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func NewRepository(db *postgres.Database) *repository {
	return &repository{
		db: db,
	}
}
