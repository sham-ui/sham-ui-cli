package category

import (
	"cms/internal/model"
	"cms/pkg/postgres"
	"cms/pkg/tracing"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const scopeName = "repository.postgres.category"

const (
	categoryNameUniqueConstraint = "category_name_unique"
	categorySlugUniqueConstraint = "category_slug_unique"
)

func extractFindOneError[T any](value T, err error) (T, error) { //nolint:ireturn
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return value, model.ErrCategoryNotFound
	case err != nil:
		return value, fmt.Errorf("query: %w", err)
	default:
		return value, nil
	}
}

type repository struct {
	db *postgres.Database
}

func (r *repository) Create(ctx context.Context, data model.Category) error {
	const op = "Create"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO category(name, slug) VALUES ($1, $2)",
		data.Name,
		data.Slug,
	)
	switch {
	case postgres.IsUniqueViolationError(err, categoryNameUniqueConstraint):
		return model.ErrCategoryNameAlreadyExists
	case postgres.IsUniqueViolationError(err, categorySlugUniqueConstraint):
		return model.ErrCategorySlugAlreadyExists
	case err != nil:
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id model.CategoryID) (*model.Category, error) {
	const op = "GetByID"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	var data category
	err := r.db.
		QueryRowContext(ctx, "SELECT id, name, slug FROM category WHERE id = $1", id).
		Scan(&data.ID, &data.Name, &data.Slug)
	if data, err = extractFindOneError(data, err); err != nil {
		return nil, err
	}
	mod := data.toModel()
	return &mod, nil
}

func (r *repository) GetBySlug(ctx context.Context, slug model.CategorySlug) (*model.Category, error) {
	const op = "GetBySlug"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	var data category
	err := r.db.
		QueryRowContext(ctx, "SELECT id, name, slug FROM category WHERE slug = $1", slug).
		Scan(&data.ID, &data.Name, &data.Slug)
	if data, err = extractFindOneError(data, err); err != nil {
		return nil, err
	}
	mod := data.toModel()
	return &mod, nil
}

func (r *repository) All(ctx context.Context) ([]model.Category, error) {
	const op = "All"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	rows, err := r.db.QueryContext(ctx, "SELECT id, name, slug FROM category ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()
	var categories []model.Category
	for rows.Next() {
		var data category
		err := rows.Scan(&data.ID, &data.Name, &data.Slug)
		if err != nil {
			return nil, fmt.Errorf("query: %w", err)
		}
		categories = append(categories, data.toModel())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return categories, nil
}

func (r *repository) Update(ctx context.Context, data model.Category) error {
	const op = "Update"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	_, err := r.db.ExecContext(
		ctx,
		"UPDATE category SET name = $2, slug = $3 WHERE id = $1",
		data.ID,
		data.Name,
		data.Slug,
	)
	switch {
	case postgres.IsUniqueViolationError(err, categoryNameUniqueConstraint):
		return model.ErrCategoryNameAlreadyExists
	case postgres.IsUniqueViolationError(err, categorySlugUniqueConstraint):
		return model.ErrCategorySlugAlreadyExists
	case err != nil:
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id model.CategoryID) error {
	const op = "Delete"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	if _, err := r.db.ExecContext(ctx, "DELETE FROM category WHERE id = $1", id); err != nil {
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func NewRepository(db *postgres.Database) *repository {
	return &repository{
		db: db,
	}
}
