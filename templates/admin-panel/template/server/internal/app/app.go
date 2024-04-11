package app

import (
	"{{ shortName }}/assets"
	"{{ shortName }}/config"
	"{{ shortName }}/internal/controller/http"
	memberRepository "{{ shortName }}/internal/repository/postgres/member"
	"{{ shortName }}/internal/service/password"
	"{{ shortName }}/internal/service/session"
	"{{ shortName }}/migrations"
	"{{ shortName }}/pkg/graceful_shutdown"
	"{{ shortName }}/pkg/postgres"
	"{{ shortName }}/pkg/tracing"
	"time"

	"github.com/go-logr/logr"
)

const (
	maxHttpSessionAge = 24 * time.Hour
)

func Run(log logr.Logger, cfg *config.Config) int {
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

	sessionSrv, err := session.New(database, cfg.Session.Secret, cfg.Session.Domain, maxHttpSessionAge)
	if err != nil {
		shutdowner.FailNotify(err, "can't create session service")
		return 1
	}

	httpServer, err := http.New(
		log,
		tracerProvider,
		propagator,
		cfg.Server,
		assetsFS,
		httpServerDependencies{
			passwordService: password.New(),
			sessionService:  sessionSrv,
			memberService:   memberRepository.NewRepository(database),
		},
	)
	if err != nil {
		shutdowner.FailNotify(err, "can't create http server")
		return 1
	}
	shutdowner.RegistryTask(httpServer)
	shutdowner.RegistryNotifier(httpServer)
	httpServer.Start()

	return 0
}
