package http

import (
	adminMembersCreateHandler "cms/internal/controller/http/handler/api/admin/members/create"
	adminMembersListHandler "cms/internal/controller/http/handler/api/admin/members/list"
	adminMembersRemoveHandler "cms/internal/controller/http/handler/api/admin/members/remove"
	adminMembersResetPasswordHandler "cms/internal/controller/http/handler/api/admin/members/reset_password"
	adminMembersUpdateHandler "cms/internal/controller/http/handler/api/admin/members/update"
	adminServerInfoHandler "cms/internal/controller/http/handler/api/admin/server_info"
	articlesArticleCreateHandler "cms/internal/controller/http/handler/api/articles/articles/create"
	articlesArticleDetailHandler "cms/internal/controller/http/handler/api/articles/articles/detail"
	articlesArticleListHandler "cms/internal/controller/http/handler/api/articles/articles/list"
	articlesArticleDeleteHandler "cms/internal/controller/http/handler/api/articles/articles/remove"
	articlesArticleUpdateHandler "cms/internal/controller/http/handler/api/articles/articles/update"
	articlesCategoryCreateHandler "cms/internal/controller/http/handler/api/articles/categories/create"
	articlesCategoryListHandler "cms/internal/controller/http/handler/api/articles/categories/list"
	articlesCategoryRemoveHandler "cms/internal/controller/http/handler/api/articles/categories/remove"
	articlesCategoryUpdateHandler "cms/internal/controller/http/handler/api/articles/categories/update"
	articlesTagCreateHandler "cms/internal/controller/http/handler/api/articles/tags/create"
	articlesTagListHandler "cms/internal/controller/http/handler/api/articles/tags/list"
	articlesTagRemoveHandler "cms/internal/controller/http/handler/api/articles/tags/remove"
	articlesTagUpdateHandler "cms/internal/controller/http/handler/api/articles/tags/update"
	assetsUploadHandler "cms/internal/controller/http/handler/api/assets/upload"
	membersEmailHandler "cms/internal/controller/http/handler/api/members/email"
	membersNameHandler "cms/internal/controller/http/handler/api/members/name"
	membersPasswordHandler "cms/internal/controller/http/handler/api/members/password"
	sessionCsrfHandler "cms/internal/controller/http/handler/api/session/csrf"
	sessonnLoginHandler "cms/internal/controller/http/handler/api/session/login"
	sessionLogoutHandler "cms/internal/controller/http/handler/api/session/logout"
	sessionValidateHandler "cms/internal/controller/http/handler/api/session/validate"
	assetsHandler "cms/internal/controller/http/handler/assets"
	fsHandler "cms/internal/controller/http/handler/fs"
	healthzHandler "cms/internal/controller/http/handler/healthz"
	metricsHandler "cms/internal/controller/http/handler/metrics"
	"cms/pkg/asserts"
	"cms/pkg/tracing"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/go-logr/logr/testr"
	"github.com/gorilla/mux"
)

func TestRouterMatch(t *testing.T) {
	// Arrange
	exporter, err := tracing.NewExporter(tracing.Config{}) //nolint:exhaustruct
	asserts.NoError(t, err)

	deps := NewMockHandlerDependencyProvider(t)
	deps.EXPECT().PasswordService().Return(nil)
	deps.EXPECT().SessionService().Return(nil)
	deps.EXPECT().SlugifyService().Return(nil)
	deps.EXPECT().MemberService().Return(nil)
	deps.EXPECT().ArticleCategoryService().Return(nil)
	deps.EXPECT().ArticleTagService().Return(nil)
	deps.EXPECT().ArticleService().Return(nil)
	deps.EXPECT().AssetsService().Return(nil)

	router, err := newRouter(
		true,
		nil,
		testr.New(t),
		tracing.NewProvider(exporter, tracing.Config{}), //nolint:exhaustruct
		tracing.NewPropagator(),
		fstest.MapFS{"index.html": {Data: []byte("<h1>Hello World</h1>"), Mode: fs.ModePerm}},
		fstest.MapFS{},
		deps,
	)
	asserts.NoError(t, err)

	testCases := []struct {
		method        string
		url           string
		expectedRoute string
	}{
		{
			method:        http.MethodGet,
			url:           "/metrics",
			expectedRoute: metricsHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/healthz",
			expectedRoute: healthzHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/api/csrftoken",
			expectedRoute: sessionCsrfHandler.RouteName,
		},
		{
			method:        http.MethodPost,
			url:           "/api/login",
			expectedRoute: sessonnLoginHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/api/validsession",
			expectedRoute: sessionValidateHandler.RouteName,
		},
		{
			method:        http.MethodPost,
			url:           "/api/logout",
			expectedRoute: sessionLogoutHandler.RouteName,
		},
		{
			method:        http.MethodPut,
			url:           "/api/members/email",
			expectedRoute: membersEmailHandler.RouteName,
		},
		{
			method:        http.MethodPut,
			url:           "/api/members/name",
			expectedRoute: membersNameHandler.RouteName,
		},
		{
			method:        http.MethodPut,
			url:           "/api/members/password",
			expectedRoute: membersPasswordHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/api/admin/members",
			expectedRoute: adminMembersListHandler.RouteName,
		},
		{
			method:        http.MethodPost,
			url:           "/api/admin/members",
			expectedRoute: adminMembersCreateHandler.RouteName,
		},
		{
			method:        http.MethodPut,
			url:           "/api/admin/members/42",
			expectedRoute: adminMembersUpdateHandler.RouteName,
		},
		{
			method:        http.MethodPut,
			url:           "/api/admin/members/42/password",
			expectedRoute: adminMembersResetPasswordHandler.RouteName,
		},
		{
			method:        http.MethodDelete,
			url:           "/api/admin/members/42",
			expectedRoute: adminMembersRemoveHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/api/admin/server-info",
			expectedRoute: adminServerInfoHandler.RouteName,
		},
		{
			method:        http.MethodPost,
			url:           "/api/categories",
			expectedRoute: articlesCategoryCreateHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/api/categories",
			expectedRoute: articlesCategoryListHandler.RouteName,
		},
		{
			method:        http.MethodPut,
			url:           "/api/categories/42",
			expectedRoute: articlesCategoryUpdateHandler.RouteName,
		},
		{
			method:        http.MethodDelete,
			url:           "/api/categories/42",
			expectedRoute: articlesCategoryRemoveHandler.RouteName,
		},
		{
			method:        http.MethodPost,
			url:           "/api/tags",
			expectedRoute: articlesTagCreateHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/api/tags",
			expectedRoute: articlesTagListHandler.RouteName,
		},
		{
			method:        http.MethodPut,
			url:           "/api/tags/42",
			expectedRoute: articlesTagUpdateHandler.RouteName,
		},
		{
			method:        http.MethodDelete,
			url:           "/api/tags/42",
			expectedRoute: articlesTagRemoveHandler.RouteName,
		},
		{
			method:        http.MethodPost,
			url:           "/api/articles",
			expectedRoute: articlesArticleCreateHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/api/articles",
			expectedRoute: articlesArticleListHandler.RouteName,
		},
		{
			method:        http.MethodPut,
			url:           "/api/articles/42",
			expectedRoute: articlesArticleUpdateHandler.RouteName,
		},
		{
			method:        http.MethodDelete,
			url:           "/api/articles/42",
			expectedRoute: articlesArticleDeleteHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/api/articles/42",
			expectedRoute: articlesArticleDetailHandler.RouteName,
		},
		{
			method:        http.MethodPost,
			url:           "/api/upload-image",
			expectedRoute: assetsUploadHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/assets/files.txt",
			expectedRoute: assetsHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/",
			expectedRoute: fsHandler.RouteName,
		},
		{
			method:        http.MethodGet,
			url:           "/page",
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
