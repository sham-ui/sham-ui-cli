package http

import (
	adminMembersCreateHandler "cms/internal/controller/http/handler/api/admin/members/create"
	adminMembersListHandler "cms/internal/controller/http/handler/api/admin/members/list"
	adminMembersDeleteHandler "cms/internal/controller/http/handler/api/admin/members/remove"
	adminMembersResetPasswordHandler "cms/internal/controller/http/handler/api/admin/members/reset_password"
	adminMembersUpdateHandler "cms/internal/controller/http/handler/api/admin/members/update"
	adminServerInfoHandler "cms/internal/controller/http/handler/api/admin/server_info"
	articlesCreateHandler "cms/internal/controller/http/handler/api/articles/articles/create"
	articlesDetailHandler "cms/internal/controller/http/handler/api/articles/articles/detail"
	articlesListHandler "cms/internal/controller/http/handler/api/articles/articles/list"
	articlesDeleteHandler "cms/internal/controller/http/handler/api/articles/articles/remove"
	articlesUpdateHandler "cms/internal/controller/http/handler/api/articles/articles/update"
	articlesCategoryCreateHandler "cms/internal/controller/http/handler/api/articles/categories/create"
	articlesCategoryListHandler "cms/internal/controller/http/handler/api/articles/categories/list"
	articlesCategoryRemoveHandler "cms/internal/controller/http/handler/api/articles/categories/remove"
	articlesCategoryUpdateHandler "cms/internal/controller/http/handler/api/articles/categories/update"
	articlesTagCreateHandler "cms/internal/controller/http/handler/api/articles/tags/create"
	articlesTagListHandler "cms/internal/controller/http/handler/api/articles/tags/list"
	articlesTagRemoveHandler "cms/internal/controller/http/handler/api/articles/tags/remove"
	articlesTagUpdateHandler "cms/internal/controller/http/handler/api/articles/tags/update"
	assetsUploadHandler "cms/internal/controller/http/handler/api/assets/upload"
	"cms/internal/controller/http/handler/api/members/email"
	"cms/internal/controller/http/handler/api/members/name"
	"cms/internal/controller/http/handler/api/members/password"
	csrfHandler "cms/internal/controller/http/handler/api/session/csrf"
	"cms/internal/controller/http/handler/api/session/login"
	"cms/internal/controller/http/handler/api/session/logout"
	"cms/internal/controller/http/handler/api/session/validate"
	assetsFileServerHandler "cms/internal/controller/http/handler/assets"
	fsHandler "cms/internal/controller/http/handler/fs"
	healthzHandler "cms/internal/controller/http/handler/healthz"
	htmlHandler "cms/internal/controller/http/handler/html"
	metricsHandler "cms/internal/controller/http/handler/metrics"
	"cms/internal/controller/http/middleware/authenticated"
	"cms/internal/controller/http/middleware/context_logger"
	corsMiddleware "cms/internal/controller/http/middleware/cors"
	"cms/internal/controller/http/middleware/csrf"
	mwLogger "cms/internal/controller/http/middleware/logger"
	"cms/internal/controller/http/middleware/recovery"
	"cms/internal/controller/http/middleware/session"
	"cms/internal/controller/http/middleware/superuser"
	"cms/internal/controller/http/middleware/tracing"
	"fmt"
	"io"
	"io/fs"
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
	contentAssets fs.FS,
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
			members.Use(
				authenticated.Middleware,
			)

			email.Setup(members, deps.SessionService(), deps.MemberService())
			name.Setup(members, deps.SessionService(), deps.MemberService())
			password.Setup(members, deps.MemberService(), deps.PasswordService())
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

		// Articles
		{
			articles := api.PathPrefix("/").Subrouter()
			articles.Use(
				authenticated.Middleware,
			)

			// Categories
			{
				categories := articles.PathPrefix("/categories").Subrouter()

				articlesCategoryCreateHandler.Setup(categories, deps.ArticleCategoryService(), deps.SlugifyService())
				articlesCategoryListHandler.Setup(categories, deps.ArticleCategoryService())
				articlesCategoryUpdateHandler.Setup(categories, deps.ArticleCategoryService(), deps.SlugifyService())
				articlesCategoryRemoveHandler.Setup(categories, deps.ArticleCategoryService())
			}

			// Tags
			{
				tags := articles.PathPrefix("/tags").Subrouter()

				articlesTagCreateHandler.Setup(tags, deps.ArticleTagService(), deps.SlugifyService())
				articlesTagListHandler.Setup(tags, deps.ArticleTagService())
				articlesTagUpdateHandler.Setup(tags, deps.ArticleTagService(), deps.SlugifyService())
				articlesTagRemoveHandler.Setup(tags, deps.ArticleTagService())
			}

			// Articles
			{
				art := articles.PathPrefix("/articles").Subrouter()

				articlesCreateHandler.Setup(art, deps.ArticleService(), deps.SlugifyService())
				articlesUpdateHandler.Setup(art, deps.ArticleService(), deps.SlugifyService())
				articlesDetailHandler.Setup(art, deps.ArticleService())
				articlesListHandler.Setup(art, deps.ArticleService())
				articlesDeleteHandler.Setup(art, deps.ArticleService())
			}
		}

		// Content assets
		{
			upload := api.PathPrefix("/").Subrouter()
			upload.Use(
				authenticated.Middleware,
			)
			assetsUploadHandler.Setup(upload, deps.AssetsService())
		}
	}

	// Content assets
	{
		assetsRouter := router.PathPrefix("/assets").Subrouter()
		assetsRouter.Use(
			authenticated.Middleware,
		)

		assetsFileServerHandler.Setup(assetsRouter, contentAssets)
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
