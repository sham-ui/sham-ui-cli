package ssr

import (
	"context"
	"net/http"
	"net/url"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name ServerSideRender --inpackage --testonly
type ServerSideRender interface {
	Render(ctx context.Context, url *url.URL, cookies []*http.Cookie) ([]byte, error)
}
