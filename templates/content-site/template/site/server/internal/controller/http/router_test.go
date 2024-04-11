package http

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	articleDetailHandler "site/internal/controller/http/handler/api/articles/detail"
	articleListCategoryHandler "site/internal/controller/http/handler/api/articles/list/category"
	articleListHandler "site/internal/controller/http/handler/api/articles/list/default"
	articleListTagHandler "site/internal/controller/http/handler/api/articles/list/tag"
	assetsHandler "site/internal/controller/http/handler/assets"
	fsHandler "site/internal/controller/http/handler/fs"
	healthzHandler "site/internal/controller/http/handler/healthz"
	metricsHandler "site/internal/controller/http/handler/metrics"
	"site/pkg/asserts"
	"site/pkg/tracing"
	"testing"
	"testing/fstest"

	"github.com/go-logr/logr/testr"
	"github.com/gorilla/mux"
)

func TestRouterMatch(t *testing.T) {
	// Arrange
	exporter, err := tracing.NewExporter(tracing.Config{}) //nolint:exhaustruct
	asserts.NoError(t, err)
	router := newRouter(
		true,
		testr.New(t),
		tracing.NewProvider(exporter, tracing.Config{}), //nolint:exhaustruct
		tracing.NewPropagator(),
		fstest.MapFS{"index.html": {Data: []byte("<h1>Hello World</h1>"), Mode: fs.ModePerm}},
		nil,
		nil,
		nil,
	)

	testCases := []struct {
		method        string
		url           string
		expectedRoute string
	}{
		{
			method:        http.MethodGet,
			url:           "/healthz",
			expectedRoute: healthzHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/metrics",
			expectedRoute: metricsHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/api/articles/article-slug",
			expectedRoute: articleDetailHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/api/articles?category=category-slug",
			expectedRoute: articleListCategoryHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/api/articles?tag=tag-slug",
			expectedRoute: articleListTagHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/api/articles",
			expectedRoute: articleListHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/assets/index.html",
			expectedRoute: assetsHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/",
			expectedRoute: fsHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/page-slug",
			expectedRoute: fsHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/contact",
			expectedRoute: fsHandler.RouteName,
		},
	}

	for _, test := range testCases {
		t.Run(test.method+test.url, func(t *testing.T) {
			// Act
			match := mux.RouteMatch{} //nolint:exhaustruct
			matched := router.Match(httptest.NewRequest(test.method, test.url, nil), &match)

			// Assert
			asserts.Equals(t, true, matched)
			asserts.Equals(t, test.expectedRoute, match.Route.GetName())
		})
	}
}
