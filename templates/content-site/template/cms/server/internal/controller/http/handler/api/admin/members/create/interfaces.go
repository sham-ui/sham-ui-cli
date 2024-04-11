package create

import (
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name MemberService --inpackage --testonly
type MemberService interface {
	Create(ctx context.Context, data model.MemberWithPassword) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name PasswordService --inpackage --testonly
type PasswordService interface {
	Hash(ctx context.Context, raw string) (model.MemberHashedPassword, error)
}
