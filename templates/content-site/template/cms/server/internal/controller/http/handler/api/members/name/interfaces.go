package name

import (
	"cms/internal/model"
	"context"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name SessionService --inpackage --testonly
type SessionService interface {
	UpdateName(w http.ResponseWriter, r *http.Request, name string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name MemberService --inpackage --testonly
type MemberService interface {
	UpdateName(ctx context.Context, id model.MemberID, name string) error
}
