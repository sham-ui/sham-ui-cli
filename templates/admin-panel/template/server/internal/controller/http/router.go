package http

import (
	"fmt"
	"io"
	"io/fs"
	adminMembersCreateHandler "{{ shortName }}/internal/controller/http/handler/api/admin/members/create"
	adminMembersListHandler "{{ shortName }}/internal/controller/http/handler/api/admin/members/list"
	adminMembersDeleteHandler "{{ shortName }}/internal/controller/http/handler/api/admin/members/remove"
	adminMembersResetPasswordHandler "{{ shortName }}/internal/controller/http/handler/api/admin/members/reset_password"
	adminMembersUpdateHandler "{{ shortName }}/internal/controller/http/handler/api/admin/members/update"
	adminServerInfoHandler "{{ shortName }}/internal/controller/http/handler/api/admin/server_info"
	"{{ shortName }}/internal/controller/http/handler/api/members/email"
	"{{ shortName }}/internal/controller/http/handler/api/members/name"
	"{{ shortName }}/internal/controller/http/handler/api/members/password"
	{{#if signupEnabled}}
	membersSignupHandler "{{ shortName }}/internal/controller/http/handler/api/members/signup"
	{{/if}}
	csrfHandler "{{ shortName }}/internal/controller/http/handler/api/session/csrf"
	"{{ shortName }}/internal/controller/http/handler/api/session/login"
	"{{ shortName }}/internal/controller/http/handler/api/session/logout"
	"{{ shortName }}/internal/controller/http/handler/api/session/validate"
	fsHandler "{{ shortName }}/internal/controller/http/handler/fs"
	healthzHandler "{{ shortName }}/internal/controller/http/handler/healthz"
	htmlHandler "{{ shortName }}/internal/controller/http/handler/html"
	metricsHandler "{{ shortName }}/internal/controller/http/handler/metrics"
	"{{ shortName }}/internal/controller/http/middleware/authenticated"
	"{{ shortName }}/internal/controller/http/middleware/context_logger"
	corsMiddleware "{{ shortName }}/internal/controller/http/middleware/cors"
	"{{ shortName }}/internal/controller/http/middleware/csrf"
	mwLogger "{{ shortName }}/internal/controller/http/middleware/logger"
	"{{ shortName }}/internal/controller/http/middleware/recovery"
	"{{ shortName }}/internal/controller/http/middleware/session"
	"{{ shortName }}/internal/controller/http/middleware/superuser"
	"{{ shortName }}/internal/controller/http/middleware/tracing"
	"time"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	csrfRequestHeader = "X-CSRF-Token"
	csrfCookieName    = "cms_csrf"
)

//nolint:funlen
func newRouter(
	cors bool,
	authKey []byte,
	logger logr.Logger,
	tracerProvider trace.TracerProvider,
	propagator propagation.TraceContext,
	cmsAssets fs.FS,
	deps HandlerDependencyProvider,
) (*mux.Router, error) {
	router := mux.NewRouter()

	// Middleware
	router.Use(
		recovery.New(logger).Middleware,
		csrf.New(authKey, csrfRequestHeader, csrfCookieName),
		tracing.New(tracerProvider, propagator),
		mwLogger.New(logger).Middleware,
		context_logger.New(logger).Middleware,
		session.New(deps.SessionService()).Middleware,
	)

	// K8s probe
	healthzHandler.Setup(router)

	// Metrics
	metricsHandler.Setup(router)

	// API
	{
		api := router.PathPrefix("/api").Subrouter()

		if cors {
			corsMiddleware.Setup(csrfRequestHeader, api)
		}

		// Session
		{
			csrfHandler.Setup(api, csrfRequestHeader)
			login.Setup(api, deps.SessionService(), deps.MemberService(), deps.PasswordService(), csrfRequestHeader)
			validate.Setup(api)
			logout.Setup(api, deps.SessionService())
		}

		// Member
		{
			members := api.PathPrefix("/members").Subrouter()

			{{#if signupEnabled}}
			// Signup
			{
				membersSignupHandler.Setup(members, deps.MemberService(), deps.PasswordService())
			}
			{{/if}}

			// Authenticated
			{
				membersAuthenticated := members.PathPrefix("/").Subrouter()
				membersAuthenticated.Use(
					authenticated.Middleware,
				)

				email.Setup(membersAuthenticated, deps.SessionService(), deps.MemberService())
				name.Setup(membersAuthenticated, deps.SessionService(), deps.MemberService())
				password.Setup(membersAuthenticated, deps.MemberService(), deps.PasswordService())
			}
		}

		// Admin
		{
			admin := api.PathPrefix("/admin").Subrouter()
			admin.Use(
				superuser.Middleware,
			)

			// Members
			{
				members := admin.PathPrefix("/members").Subrouter()

				adminMembersListHandler.Setup(members, deps.MemberService())
				adminMembersCreateHandler.Setup(members, deps.MemberService(), deps.PasswordService())
				adminMembersUpdateHandler.Setup(members, deps.MemberService())
				adminMembersResetPasswordHandler.Setup(members, deps.MemberService(), deps.PasswordService())
				adminMembersDeleteHandler.Setup(members, deps.MemberService())
			}
			adminServerInfoHandler.Setup(admin, cmsAssets, time.Now())
		}
	}

	// Spa
	{
		spaRouter := router.PathPrefix("/").Subrouter().StrictSlash(true)

		htmlFile, err := cmsAssets.Open("index.html")
		if err != nil {
			return nil, fmt.Errorf("failed to open index.html: %w", err)
		}
		htmlContent, err := io.ReadAll(htmlFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read index.html: %w", err)
		}
		html := htmlHandler.New(htmlContent)
		fsHandler.Setup(spaRouter, cmsAssets, html)
	}

	return router, nil
}
