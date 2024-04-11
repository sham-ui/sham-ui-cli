package remove

import (
	"cms/internal/model"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name TagService --inpackage --testonly
type TagService interface {
	Delete(ctx context.Context, id model.TagID) error
}
