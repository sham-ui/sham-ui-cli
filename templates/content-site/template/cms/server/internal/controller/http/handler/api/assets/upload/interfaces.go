package upload

import (
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name AssetsService --inpackage --testonly
type AssetsService interface {
	Save(ctx context.Context, filename string, content []byte) (*model.Asset, error)
}
