package asset

import "context"

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name AssetService --inpackage --testonly
type AssetService interface {
	ReadFile(ctx context.Context, path string) ([]byte, error)
}
