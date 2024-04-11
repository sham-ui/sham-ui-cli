package password

import (
	"context"
	"{{ shortName }}/internal/model"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name MemberService --inpackage --testonly
type MemberService interface {
	UpdatePassword(ctx context.Context, id model.MemberID, password model.MemberHashedPassword) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name PasswordService --inpackage --testonly
type PasswordService interface {
	Hash(ctx context.Context, raw string) (model.MemberHashedPassword, error)
}
