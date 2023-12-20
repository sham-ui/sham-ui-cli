package main

import (
	"os"
	"site/assets"
	"site/config"
	"site/internal/controller/http"
	"site/internal/external_api/cms"
	"site/internal/service/ssr"
	"site/pkg/graceful_shutdown"
	"site/pkg/logger"
	"site/pkg/tracing"
	ssrFiles "site/ssr"
)

func app() int { //nolint:funlen
	logr := logger.NewLogger(0)

	cfg, err := config.LoadConfiguration(logr, "config.cfg")
	if err != nil {
		logr.Error(err, "can't load config")
		os.Exit(1)
	}
	logr.Info("config loaded", "config", cfg)

	shutdowner := graceful_shutdown.New(logr)
	defer shutdowner.Wait()

	spanExporter, err := tracing.NewExporter(cfg.Tracer)
	if err != nil {
		logr.Error(err, "can't create span exporter")
		return 1
	}
	shutdowner.RegistryTask(spanExporter)
	tracerProvider := tracing.NewProvider(spanExporter, cfg.Tracer)
	propagator := tracing.NewPropagator()

	apiClient, err := cms.New(tracerProvider, propagator, cfg.API)
	if err != nil {
		logr.Error(err, "can't create cms client")
		return 1
	}
	logr.Info("cms client created")

	files, err := assets.Files()
	if err != nil {
		logr.Error(err, "can't load files")
		return 1
	}

	ssrSrv := ssr.New(
		logr,
		tracerProvider,
		propagator,
		cfg.Server.URL()+http.APIPrefix,
		ssrFiles.Script,
	)
	if err := ssrSrv.Start(); err != nil {
		logr.Error(err, "can't start ssr server")
		return 1
	}
	shutdowner.RegistryTask(ssrSrv)

	httpServer := http.New(
		logr,
		tracerProvider,
		propagator,
		cfg.Server,
		files,
		apiClient,
		apiClient,
		ssrSrv,
	)
	shutdowner.RegistryTask(httpServer)
	shutdowner.RegistryNotifier(httpServer)
	httpServer.Start()

	return 0
}

func main() {
	os.Exit(app())
}
