package list

import (
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name TagService --inpackage --testonly
type TagService interface {
	All(ctx context.Context) ([]model.Tag, error)
}
