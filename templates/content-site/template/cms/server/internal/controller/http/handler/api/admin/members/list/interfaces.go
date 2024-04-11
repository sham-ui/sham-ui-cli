package list

import (
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name MemberService --inpackage --testonly
type MemberService interface {
	Find(ctx context.Context, offset, limit int64) ([]model.Member, error)
	Total(ctx context.Context) (int, error)
}
