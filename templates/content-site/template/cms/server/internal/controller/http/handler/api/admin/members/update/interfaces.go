package update

import (
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name MemberService --inpackage --testonly
type MemberService interface {
	Update(ctx context.Context, data model.Member) error
}
