package assets

import (
	"context"

	"site/internal/model"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name Service --inpackage --testonly
type Service interface {
	Asset(ctx context.Context, path string) (*model.Asset, error)
}
