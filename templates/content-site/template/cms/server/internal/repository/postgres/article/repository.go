package article

import (
	"cms/internal/model"
	"cms/pkg/postgres"
	"cms/pkg/tracing"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const scopeName = "repository.postgres.article"

const (
	articleNameUniqueConstraint     = "article_name_unique"
	articleSlugUniqueConstraint     = "article_slug_unique"
	articleCategoryIDFKeyConstraint = "article_category_id_fkey"
)

func extractFindOneError[T any](x T, err error) (T, error) { //nolint:ireturn
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return x, model.ErrArticleNotFound
	default:
		return x, err
	}
}

type repository struct {
	db *postgres.Database
}

func (r *repository) Create(ctx context.Context, data model.Article) (model.ArticleID, error) {
	const op = "Create"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	var id string
	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO article"+
			"(title, slug, category_id, short_body, body, published_at) "+
			"VALUES ($1, $2, $3, $4, $5, $6) "+
			"RETURNING id",
		data.Title,
		data.Slug,
		data.CategoryID,
		data.ShortBody,
		data.Body,
		data.PublishedAt,
	).Scan(&id)
	switch {
	case postgres.IsUniqueViolationError(err, articleNameUniqueConstraint):
		return "", model.ErrArticleTitleAlreadyExists
	case postgres.IsUniqueViolationError(err, articleSlugUniqueConstraint):
		return "", model.ErrArticleSlugAlreadyExists
	case postgres.IsForeignKeyViolationError(err, articleCategoryIDFKeyConstraint):
		return "", model.ErrCategoryNotFound
	case err != nil:
		return "", fmt.Errorf("query: %w", err)
	}
	return model.ArticleID(id), nil
}

func (r *repository) FindByID(ctx context.Context, id model.ArticleID) (model.Article, error) {
	const op = "FindByID"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	var item article
	err := r.db.QueryRowContext(
		ctx,
		"SELECT "+
			"id, title, slug, category_id, short_body, body, published_at "+
			"FROM article "+
			"WHERE id = $1",
		id,
	).Scan(
		&item.ID,
		&item.Title,
		&item.Slug,
		&item.CategoryID,
		&item.ShortBody,
		&item.Body,
		&item.PublishedAt,
	)
	if item, err = extractFindOneError(item, err); err != nil {
		return model.Article{}, err
	}
	return item.toModel(), nil
}

func (r *repository) FindBySlug(ctx context.Context, slug model.ArticleSlug) (model.Article, error) {
	const op = "FindBySlug"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	var item article
	err := r.db.QueryRowContext(
		ctx,
		"SELECT "+
			"id, title, slug, category_id, short_body, body, published_at "+
			"FROM article "+
			"WHERE slug = $1",
		slug,
	).Scan(
		&item.ID,
		&item.Title,
		&item.Slug,
		&item.CategoryID,
		&item.ShortBody,
		&item.Body,
		&item.PublishedAt,
	)
	if item, err = extractFindOneError(item, err); err != nil {
		return model.Article{}, err
	}
	return item.toModel(), nil
}

func (r *repository) FindShortInfo(ctx context.Context, offset, limit int64) ([]model.ArticleShortInfo, error) {
	const op = "FindShortInfo"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT "+
			"id, title, slug, category_id, published_at "+
			"FROM article "+
			"ORDER BY published_at DESC "+
			"LIMIT $1 "+
			"OFFSET $2",

		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var articles []model.ArticleShortInfo
	for rows.Next() {
		var item article
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Slug,
			&item.CategoryID,
			&item.PublishedAt,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		articles = append(articles, item.toModelShortInfo())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return articles, nil
}

func (r *repository) FindShortInfoWithCategory(
	ctx context.Context,
	offset, limit int64,
) ([]model.ArticleShortInfoWithCategory, error) {
	const op = "FindShortInfoWithCategory"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT "+
			"a.id, a.title, a.slug, a.short_body, ct.id, ct.name, ct.slug, a.published_at "+
			"FROM article a "+
			"JOIN category ct ON ct.id = a.category_id "+
			"ORDER BY published_at DESC "+
			"LIMIT $1 "+
			"OFFSET $2",
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var articles []model.ArticleShortInfoWithCategory
	for rows.Next() {
		var item article
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Slug,
			&item.ShortBody,
			&item.CategoryID,
			&item.CategoryName,
			&item.CategorySlug,
			&item.PublishedAt,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		articles = append(articles, item.toModelWithCategory())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return articles, nil
}

func (r *repository) FindShortInfoWithCategoryForTag(
	ctx context.Context,
	tagSlug model.TagSlug,
	offset, limit int64,
) ([]model.ArticleShortInfoWithCategory, error) {
	const op = "FindShortInfoWithCategoryForTag"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT "+
			"a.id, a.title, a.slug, a.short_body, ct.id, ct.name, ct.slug, a.published_at "+
			"FROM article a "+
			"JOIN category ct ON ct.id = a.category_id "+
			"JOIN article_tag a_t ON a_t.article_id = a.id "+
			"JOIN tag t ON t.id = a_t.tag_id "+
			"WHERE t.slug = $1 "+
			"ORDER BY published_at DESC "+
			"LIMIT $2 "+
			"OFFSET $3",
		tagSlug,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()
	var articles []model.ArticleShortInfoWithCategory
	for rows.Next() {
		var item article
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Slug,
			&item.ShortBody,
			&item.CategoryID,
			&item.CategoryName,
			&item.CategorySlug,
			&item.PublishedAt,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		articles = append(articles, item.toModelWithCategory())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return articles, nil
}

func (r *repository) FindShortInfoWithCategoryForCategory(
	ctx context.Context,
	categorySlug model.CategorySlug,
	offset, limit int64,
) ([]model.ArticleShortInfoWithCategory, error) {
	const op = "FindShortInfoWithCategoryForCategory"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT "+
			"a.id, a.title, a.slug, a.short_body, ct.id, ct.name, ct.slug, a.published_at "+
			"FROM article a "+
			"JOIN category ct ON ct.id = a.category_id "+
			"WHERE ct.slug = $1 "+
			"ORDER BY published_at DESC "+
			"LIMIT $2 "+
			"OFFSET $3",
		categorySlug,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()
	var articles []model.ArticleShortInfoWithCategory
	for rows.Next() {
		var item article
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Slug,
			&item.ShortBody,
			&item.CategoryID,
			&item.CategoryName,
			&item.CategorySlug,
			&item.PublishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		articles = append(articles, item.toModelWithCategory())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return articles, nil
}

func (r *repository) FindShortInfoWithCategoryForQuery(
	ctx context.Context,
	query string,
	offset, limit int64,
) ([]model.ArticleShortInfoWithCategory, error) {
	const op = "FindShortInfoWithCategoryForQuery"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT 
    		a.id, a.title, a.slug, a.short_body, ct.id, ct.name, ct.slug, a.published_at
		FROM article a
        JOIN category ct ON ct.id = a.category_id
		WHERE a.title ILIKE $1
   			OR a.short_body ILIKE $1
		    OR a.body ILIKE $1
		ORDER BY published_at DESC
		OFFSET $3 LIMIT $2`,
		"%"+query+"%",
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()
	var articles []model.ArticleShortInfoWithCategory
	for rows.Next() {
		var item article
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Slug,
			&item.ShortBody,
			&item.CategoryID,
			&item.CategoryName,
			&item.CategorySlug,
			&item.PublishedAt,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		articles = append(articles, item.toModelWithCategory())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return articles, nil
}

func (r *repository) Total(ctx context.Context) (int, error) {
	const op = "Total"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	var count int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM article").Scan(&count); err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}
	return count, nil
}

func (r *repository) TotalForCategory(ctx context.Context, slug model.CategorySlug) (int, error) {
	const op = "TotalForCategory"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	var count int
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) 
		FROM article a 
		JOIN category ct ON ct.id = a.category_id 
		WHERE ct.slug = $1`,
		slug,
	).Scan(&count); err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}
	return count, nil
}

func (r *repository) TotalForTag(ctx context.Context, slug model.TagSlug) (int, error) {
	const op = "TotalForTag"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	var count int
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) 
		FROM article a 
		JOIN article_tag a_t ON a_t.article_id = a.id 
		JOIN tag t ON t.id = a_t.tag_id 
		WHERE t.slug = $1`,
		slug,
	).Scan(&count); err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}
	return count, nil
}

func (r *repository) TotalForQuery(ctx context.Context, query string) (int, error) {
	const op = "TotalForQuery"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	var count int
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) 
		FROM article 
		WHERE 
		    title ILIKE $1 OR 
		    short_body ILIKE $1 OR 
		    body ILIKE $1`,
		"%"+query+"%",
	).Scan(&count); err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}
	return count, nil
}

func (r *repository) Update(ctx context.Context, data model.Article) error {
	const op = "Update"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	_, err := r.db.ExecContext(
		ctx,
		"UPDATE article SET "+
			"title = $1, slug = $2, category_id = $3, short_body = $4, body = $5, published_at = $6 "+
			"WHERE id = $7 ",
		data.Title,
		data.Slug,
		data.CategoryID,
		data.ShortBody,
		data.Body,
		data.PublishedAt,
		data.ID,
	)
	switch {
	case postgres.IsUniqueViolationError(err, articleNameUniqueConstraint):
		return model.ErrArticleTitleAlreadyExists
	case postgres.IsUniqueViolationError(err, articleSlugUniqueConstraint):
		return model.ErrArticleSlugAlreadyExists
	case postgres.IsForeignKeyViolationError(err, articleCategoryIDFKeyConstraint):
		return model.ErrCategoryNotFound
	case err != nil:
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id model.ArticleID) error {
	const op = "Delete"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	if _, err := r.db.ExecContext(ctx, "DELETE FROM article WHERE id = $1", id); err != nil {
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func NewRepository(db *postgres.Database) *repository {
	return &repository{
		db: db,
	}
}
