package remove

import (
	"context"
	"{{ shortName }}/internal/model"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name MemberService --inpackage --testonly
type MemberService interface {
	Delete(ctx context.Context, id model.MemberID) error
}
