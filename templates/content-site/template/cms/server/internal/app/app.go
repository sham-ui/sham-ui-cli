package app

import (
	"cms/assets"
	"cms/config"
	"cms/internal/controller/grpc"
	"cms/internal/controller/http"
	articleRepository "cms/internal/repository/postgres/article"
	articleTagRepository "cms/internal/repository/postgres/article_tag"
	categoryRepository "cms/internal/repository/postgres/category"
	memberRepository "cms/internal/repository/postgres/member"
	tagRepository "cms/internal/repository/postgres/tag"
	"cms/internal/service/article"
	"cms/internal/service/asset"
	"cms/internal/service/password"
	"cms/internal/service/session"
	"cms/internal/service/slugify"
	"cms/migrations"
	"cms/pkg/graceful_shutdown"
	"cms/pkg/postgres"
	"cms/pkg/tracing"
	"os"
	"time"

	"github.com/go-logr/logr"
)

const (
	maxHttpSessionAge = 24 * time.Hour
)

func Run(log logr.Logger, cfg *config.Config) int { //nolint:funlen
	shutdowner := graceful_shutdown.New(log)
	defer shutdowner.Wait()

	spanExporter, err := tracing.NewExporter(cfg.Tracer)
	if err != nil {
		shutdowner.FailNotify(err, "can't create span exporter")
		return 1
	}
	shutdowner.RegistryTask(spanExporter)
	tracerProvider := tracing.NewProvider(spanExporter, cfg.Tracer)
	propagator := tracing.NewPropagator()

	assetsFS, err := assets.Files()
	if err != nil {
		shutdowner.FailNotify(err, "can't create assets")
		return 1
	}

	database, err := postgres.New(log, cfg.Database.URL())
	if err != nil {
		shutdowner.FailNotify(err, "can't connect to database")
		return 1
	}
	shutdowner.RegistryTask(database)

	if err := migrations.Up(log, database); err != nil {
		shutdowner.FailNotify(err, "can't up migrations")
		return 1
	}

	tagRepo := tagRepository.NewRepository(database)
	categoryRepo := categoryRepository.NewRepository(database)
	articleRepo := articleRepository.NewRepository(database)
	articleTagRepo := articleTagRepository.NewRepository(database)
	articleSrv := article.New(database, articleRepo, tagRepo, articleTagRepo)

	sessionSrv, err := session.New(database, cfg.Session.Secret, cfg.Session.Domain, maxHttpSessionAge)
	if err != nil {
		shutdowner.FailNotify(err, "can't create session service")
		return 1
	}

	assetSrv := asset.New(cfg.Upload.Path)

	httpServer, err := http.New(
		log,
		tracerProvider,
		propagator,
		cfg.Server,
		assetsFS,
		os.DirFS(cfg.Upload.Path),
		httpServerDependencies{
			passwordService:        password.New(),
			sessionService:         sessionSrv,
			slugifyService:         slugify.New(),
			memberService:          memberRepository.NewRepository(database),
			articleCategoryService: categoryRepo,
			articleTagService:      tagRepo,
			articleService:         articleSrv,
			assetsService:          assetSrv,
		},
	)
	if err != nil {
		shutdowner.FailNotify(err, "can't create http server")
		return 1
	}
	shutdowner.RegistryTask(httpServer)
	shutdowner.RegistryNotifier(httpServer)
	httpServer.Start()

	grpsServer, err := grpc.New(
		log,
		tracerProvider,
		propagator,
		cfg.API,
		grpcServerDependencies{
			articleService:    articleRepo,
			categoryService:   categoryRepo,
			tagService:        tagRepo,
			articleTagService: articleTagRepo,
			assetService:      assetSrv,
		},
	)
	if err != nil {
		shutdowner.FailNotify(err, "can't create grpc server")
		return 1
	}
	shutdowner.RegistryTask(grpsServer)
	shutdowner.RegistryNotifier(grpsServer)
	grpsServer.Start()

	return 0
}
