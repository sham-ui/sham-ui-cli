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

	"{{ shortName }}/config"

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
	s.srv = nil
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
	s.logger.Info("http server started", "url", s.url)
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
	deps HandlerDependencyProvider,
) (*server, error) {
	router, err := newRouter(
		cfg.CORS,
		cfg.CSRF,
		logger,
		tracerProvider,
		propagator,
		assetFS,
		deps,
	)
	if err != nil {
		return nil, err
	}
	return &server{
		srv: &http.Server{
			Addr:              cfg.Address(),
			ReadHeaderTimeout: requestHeaderTimeout,
			Handler:           router,
		},
		logger: logger,
		url:    cfg.URL(),
		lock:   sync.Mutex{},
		notify: make(chan error, 1),
	}, nil
}
