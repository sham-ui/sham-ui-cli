package tag

import (
	"cms/internal/model"
	"cms/pkg/postgres"
	"cms/pkg/tracing"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const scopeName = "repository.postgres.tag"

const (
	tagNameUniqueConstraint = "tag_name_unique"
	tagSlugUniqueConstraint = "tag_slug_unique"
)

type repository struct {
	db *postgres.Database
}

func (r *repository) Create(ctx context.Context, data model.Tag) (model.TagID, error) {
	const op = "Create"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	var id string
	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO tag(name, slug) VALUES ($1, $2) RETURNING id",
		data.Name,
		data.Slug,
	).Scan(&id)
	switch {
	case postgres.IsUniqueViolationError(err, tagNameUniqueConstraint):
		return "", model.ErrTagNameAlreadyExists
	case postgres.IsUniqueViolationError(err, tagSlugUniqueConstraint):
		return "", model.ErrTagSlugAlreadyExists
	case err != nil:
		return "", fmt.Errorf("query: %w", err)
	}
	return model.TagID(id), nil
}

func (r *repository) GetBySlug(ctx context.Context, slug model.TagSlug) (model.Tag, error) {
	const op = "GetBySlug"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	var data tag
	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, name FROM tag WHERE slug = $1",
		slug,
	).Scan(&data.ID, &data.Name)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return model.Tag{}, model.ErrTagNotFound
	case err != nil:
		return model.Tag{}, fmt.Errorf("query: %w", err)
	}
	return model.Tag{
		ID:   model.TagID(data.ID),
		Slug: slug,
		Name: data.Name,
	}, nil
}

func (r *repository) All(ctx context.Context) ([]model.Tag, error) {
	const op = "All"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	rows, err := r.db.QueryContext(ctx, "SELECT id, name, slug FROM tag ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()
	var tags []model.Tag
	for rows.Next() {
		var data tag
		err := rows.Scan(&data.ID, &data.Name, &data.Slug)
		if err != nil {
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

func (r *repository) Update(ctx context.Context, data model.Tag) error {
	const op = "Update"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	_, err := r.db.ExecContext(
		ctx,
		"UPDATE tag SET name = $2, slug = $3 WHERE id = $1",
		data.ID,
		data.Name,
		data.Slug,
	)
	switch {
	case postgres.IsUniqueViolationError(err, tagNameUniqueConstraint):
		return model.ErrTagNameAlreadyExists
	case postgres.IsUniqueViolationError(err, tagSlugUniqueConstraint):
		return model.ErrTagSlugAlreadyExists
	case err != nil:
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id model.TagID) error {
	const op = "Delete"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	if _, err := r.db.ExecContext(ctx, "DELETE FROM tag WHERE id = $1", id); err != nil {
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func NewRepository(db *postgres.Database) *repository {
	return &repository{
		db: db,
	}
}
