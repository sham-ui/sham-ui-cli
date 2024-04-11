package login

import (
	"cms/internal/model"
	"context"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name PasswordService --inpackage --testonly
type PasswordService interface {
	Compare(ctx context.Context, hashed model.MemberHashedPassword, raw string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name SessionService --inpackage --testonly
type SessionService interface {
	Create(w http.ResponseWriter, r *http.Request, member *model.Member) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name MemberService --inpackage --testonly
type MemberService interface {
	GetByEmail(ctx context.Context, email string) (*model.MemberWithPassword, error)
}
