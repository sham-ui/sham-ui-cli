package email

import (
	"cms/internal/model"
	"context"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name SessionService --inpackage --testonly
type SessionService interface {
	UpdateEmail(w http.ResponseWriter, r *http.Request, email string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name MemberService --inpackage --testonly
type MemberService interface {
	UpdateEmail(ctx context.Context, id model.MemberID, email string) error
}
