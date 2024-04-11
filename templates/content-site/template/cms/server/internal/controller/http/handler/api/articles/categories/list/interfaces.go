package list

import (
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name CategoryService --inpackage --testonly
type CategoryService interface {
	All(ctx context.Context) ([]model.Category, error)
}
