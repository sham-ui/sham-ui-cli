package http

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"site/config"

	"github.com/go-logr/logr"
)

const requestHeaderTimeout = 30 * time.Second

type server struct {
	srv    *http.Server
	logger logr.Logger
	url    string
	lock   sync.Mutex
	notify chan error
}

func (s *server) String() string {
	return "http-server(" + s.url + ")"
}

func (s *server) GracefulShutdown(ctx context.Context) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.srv == nil {
		return nil
	}
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}
	return nil
}

func (s *server) Notify() <-chan error {
	return s.notify
}

func (s *server) Start() {
	s.lock.Lock()
	srv := s.srv
	s.lock.Unlock()
	if srv == nil {
		return
	}
	s.logger.Info("server started", "url", s.url)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.notify <- fmt.Errorf("listen and serve: %w", err)
		}
	}()
}

func New(
	logger logr.Logger,
	tracerProvider trace.TracerProvider,
	propagator propagation.TraceContext,
	cfg config.Server,
	assetFS fs.FS,
	assetsService AssetsService,
	articleService ArticlesService,
	serverSideRender ServerSideRender,
) *server {
	return &server{
		srv: &http.Server{
			Addr:              cfg.Address(),
			ReadHeaderTimeout: requestHeaderTimeout,
			Handler: newRouter(
				cfg.Cors,
				logger,
				tracerProvider,
				propagator,
				assetFS,
				assetsService,
				articleService,
				serverSideRender,
			),
		},
		logger: logger,
		url:    cfg.URL(),
		lock:   sync.Mutex{},
		notify: make(chan error, 1),
	}
}
