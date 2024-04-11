package http

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	adminMembersCreateHandler "{{ shortName }}/internal/controller/http/handler/api/admin/members/create"
	adminMembersListHandler "{{ shortName }}/internal/controller/http/handler/api/admin/members/list"
	adminMembersRemoveHandler "{{ shortName }}/internal/controller/http/handler/api/admin/members/remove"
	adminMembersResetPasswordHandler "{{ shortName }}/internal/controller/http/handler/api/admin/members/reset_password"
	adminMembersUpdateHandler "{{ shortName }}/internal/controller/http/handler/api/admin/members/update"
	adminServerInfoHandler "{{ shortName }}/internal/controller/http/handler/api/admin/server_info"
	membersEmailHandler "{{ shortName }}/internal/controller/http/handler/api/members/email"
	membersNameHandler "{{ shortName }}/internal/controller/http/handler/api/members/name"
	membersPasswordHandler "{{ shortName }}/internal/controller/http/handler/api/members/password"
	{{#if signupEnabled}}
	membersSignupHandler "{{ shortName }}/internal/controller/http/handler/api/members/signup"
	{{/if}}
	sessionCsrfHandler "{{ shortName }}/internal/controller/http/handler/api/session/csrf"
	sessionLoginHandler "{{ shortName }}/internal/controller/http/handler/api/session/login"
	sessionLogoutHandler "{{ shortName }}/internal/controller/http/handler/api/session/logout"
	sessionValidateHandler "{{ shortName }}/internal/controller/http/handler/api/session/validate"
	fsHandler "{{ shortName }}/internal/controller/http/handler/fs"
	healthzHandler "{{ shortName }}/internal/controller/http/handler/healthz"
	metricsHandler "{{ shortName }}/internal/controller/http/handler/metrics"
	"{{ shortName }}/pkg/asserts"
	"{{ shortName }}/pkg/tracing"
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
	deps.EXPECT().MemberService().Return(nil)

	router, err := newRouter(
		true,
		nil,
		testr.New(t),
		tracing.NewProvider(exporter, tracing.Config{}), //nolint:exhaustruct
		tracing.NewPropagator(),
		fstest.MapFS{"index.html": {Data: []byte("<h1>Hello World</h1>"), Mode: fs.ModePerm}},
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
			expectedRoute: sessionLoginHandler.RouteName,
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
		{{#if signupEnabled}}
		{
			method:        http.MethodPost,
			url:           "/api/members",
			expectedRoute: membersSignupHandler.RouteName,
		},
		{{/if}}
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
