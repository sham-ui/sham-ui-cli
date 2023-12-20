package http

import (
	"io/fs"
	articleDetailHandler "site/internal/controller/http/handler/api/articles/detail"
	articleListCategoryHandler "site/internal/controller/http/handler/api/articles/list/category"
	articleListHandler "site/internal/controller/http/handler/api/articles/list/default"
	articleListQueryHandler "site/internal/controller/http/handler/api/articles/list/query"
	articleListTagHandler "site/internal/controller/http/handler/api/articles/list/tag"
	assetsHandler "site/internal/controller/http/handler/assets"
	fsHandler "site/internal/controller/http/handler/fs"
	healthzHandler "site/internal/controller/http/handler/healthz"
	metricsHandler "site/internal/controller/http/handler/metrics"
	ssrHandler "site/internal/controller/http/handler/ssr"
	"site/internal/controller/http/middleware/context_logger"
	"site/internal/controller/http/middleware/cors"
	mwLogger "site/internal/controller/http/middleware/logger"
	"site/internal/controller/http/middleware/recovery"
	"site/internal/controller/http/middleware/tracing"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const APIPrefix = "/api"

func newRouter(
	logger logr.Logger,
	tracerProvider trace.TracerProvider,
	propagator propagation.TraceContext,
	assetFS fs.FS,
	assetsService AssetsService,
	articleService ArticlesService,
	serverSideRender ServerSideRender,
) *mux.Router {
	router := mux.NewRouter()

	// Middleware
	router.Use(
		recovery.New(logger).Middleware,
		tracing.New(tracerProvider, propagator),
		mwLogger.New(logger).Middleware,
		context_logger.New(logger).Middleware,
	)

	// K8s probe
	healthzHandler.Setup(router)

	// Metrics
	metricsHandler.Setup(router)

	// API
	{
		api := router.PathPrefix(APIPrefix).Subrouter()
		cors.Setup(api)

		// Articles
		{
			articleDetailHandler.Setup(api, articleService)
			articleListQueryHandler.Setup(api, articleService)
			articleListCategoryHandler.Setup(api, articleService)
			articleListTagHandler.Setup(api, articleService)
			articleListHandler.Setup(api, articleService)
		}
	}

	// Content assets
	assetsHandler.Setup(router, assetsService)

	// Spa
	{
		spaRouter := router.PathPrefix("/").Subrouter().StrictSlash(true)
		ssr := ssrHandler.NewHandler(serverSideRender)
		fsHandler.Setup(spaRouter, assetFS, ssr)
	}

	return router
}
