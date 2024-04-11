package migrations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"{{ shortName }}/internal/model"
	"{{ shortName }}/internal/repository/postgres/member"
	"{{ shortName }}/internal/service/password"
	"{{ shortName }}/pkg/logger"
	"{{ shortName }}/pkg/postgres"

	"github.com/pressly/goose/v3"
)

//nolint:gochecknoinits
func init() {
	goose.AddMigrationNoTxContext(upCreateSuperuser, downCreateSuperuser)
}

var (
	//nolint:exhaustruct,gochecknoglobals
	superuserMember = model.Member{
		Email:       "root",
		Name:        "Superuser",
		IsSuperuser: true,
	}
	//nolint:gochecknoglobals
	superUserPassword = "password"
)

func upCreateSuperuser(ctx context.Context, db *sql.DB) error {
	log := logger.NewLogger(0)
	pg := postgres.NewFromConnection(log, db)
	pass, err := password.New().Hash(ctx, superUserPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	if err := member.NewRepository(pg).Create(ctx, model.MemberWithPassword{
		Member:         superuserMember,
		HashedPassword: pass,
	}); err != nil {
		return fmt.Errorf("create member: %w", err)
	}
	return nil
}

func downCreateSuperuser(ctx context.Context, db *sql.DB) error {
	log := logger.NewLogger(0)
	pg := postgres.NewFromConnection(log, db)
	repo := member.NewRepository(pg)
	mem, err := repo.GetByEmail(ctx, superuserMember.Email)
	switch {
	case errors.Is(err, model.ErrMemberNotFound):
		return nil
	case err != nil:
		return fmt.Errorf("get by email: %w", err)
	}
	if err := repo.Delete(ctx, mem.ID); err != nil {
		return fmt.Errorf("delete member: %w", err)
	}
	return nil
}
