package member

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"{{ shortName }}/internal/model"
	"{{ shortName }}/pkg/postgres"
	"{{ shortName }}/pkg/tracing"
)

const scopeName = "repository.postgres.member"

const (
	memberEmailUniqueConstraint = "member_email_unique"
)

type repository struct {
	db *postgres.Database
}

func (r *repository) Create(ctx context.Context, data model.MemberWithPassword) error {
	const op = "Create"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO members(name, email, password, is_superuser) VALUES ($1, $2, $3, $4)",
		data.Name,
		data.Email,
		data.HashedPassword,
		data.IsSuperuser,
	)
	switch {
	case postgres.IsUniqueViolationError(err, memberEmailUniqueConstraint):
		return model.ErrMemberEmailAlreadyExists
	case err != nil:
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*model.MemberWithPassword, error) {
	const op = "GetByEmail"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	var data member
	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, email, name, password, is_superuser FROM members WHERE email = $1",
		email,
	).Scan(
		&data.ID,
		&data.Email,
		&data.Name,
		&data.Password,
		&data.IsSuperuser,
	)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, model.ErrMemberNotFound
	case err != nil:
		return nil, fmt.Errorf("query: %w", err)
	}
	return &model.MemberWithPassword{
		Member: model.Member{
			ID:          model.MemberID(data.ID),
			Email:       data.Email,
			Name:        data.Name,
			IsSuperuser: data.IsSuperuser,
		},
		HashedPassword: model.MemberHashedPassword(data.Password),
	}, nil
}

func (r *repository) Find(ctx context.Context, offset, limit int64) ([]model.Member, error) {
	const op = "Find"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, email, name, is_superuser FROM members ORDER BY id LIMIT $1 OFFSET $2",
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var members []model.Member
	for rows.Next() {
		var data member
		err := rows.Scan(
			&data.ID,
			&data.Email,
			&data.Name,
			&data.IsSuperuser,
		)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		members = append(members, model.Member{
			ID:          model.MemberID(data.ID),
			Email:       data.Email,
			Name:        data.Name,
			IsSuperuser: data.IsSuperuser,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return members, nil
}

func (r *repository) Total(ctx context.Context) (int, error) {
	const op = "Total"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	var count int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM members").Scan(&count); err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}
	return count, nil
}

func (r *repository) Update(ctx context.Context, data model.Member) error {
	const op = "Update"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	_, err := r.db.ExecContext(
		ctx,
		"UPDATE members SET name = $2, email = $3, is_superuser = $4 WHERE id = $1",
		data.ID,
		data.Name,
		data.Email,
		data.IsSuperuser,
	)
	switch {
	case postgres.IsUniqueViolationError(err, memberEmailUniqueConstraint):
		return model.ErrMemberEmailAlreadyExists
	case err != nil:
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func (r *repository) UpdateEmail(ctx context.Context, id model.MemberID, newEmail string) error {
	const op = "UpdateEmail"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	_, err := r.db.ExecContext(ctx, "UPDATE members SET email = $1 WHERE id = $2", newEmail, id)
	switch {
	case postgres.IsUniqueViolationError(err, memberEmailUniqueConstraint):
		return model.ErrMemberEmailAlreadyExists
	case err != nil:
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func (r *repository) UpdateName(ctx context.Context, id model.MemberID, newName string) error {
	const op = "UpdateName"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	if _, err := r.db.ExecContext(ctx, "UPDATE members SET name = $1 WHERE id = $2", newName, id); err != nil {
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func (r *repository) UpdatePassword(ctx context.Context, id model.MemberID, password model.MemberHashedPassword) error {
	const op = "UpdatePassword"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	if _, err := r.db.ExecContext(
		ctx,
		"UPDATE members SET password = $1 WHERE id = $2",
		password,
		id,
	); err != nil {
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id model.MemberID) error {
	const op = "Delete"

	ctx, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	if _, err := r.db.ExecContext(ctx, "DELETE FROM members WHERE id = $1", id); err != nil {
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func NewRepository(db *postgres.Database) *repository {
	return &repository{
		db: db,
	}
}
