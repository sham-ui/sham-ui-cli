package grpc

import (
	"cms/config"
	contextLoggerInterceptor "cms/internal/controller/grpc/interceptor/context_logger"
	loggerInterceptor "cms/internal/controller/grpc/interceptor/logger"
	"cms/internal/controller/grpc/proto"
	"cms/pkg/net_addr"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"net"
	"os"
	"path"
	"sync"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	srv    *grpc.Server
	lis    net.Listener
	logger logr.Logger
	addr   string
	lock   sync.Mutex
	notify chan error
}

func (s *server) String() string {
	return "grpc-server(" + s.addr + ")"
}

func (s *server) GracefulShutdown(_ context.Context) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.srv == nil {
		return nil
	}
	s.srv.GracefulStop()
	s.srv = nil
	s.lis = nil
	return nil
}

func (s *server) Notify() <-chan error {
	return s.notify
}

func (s *server) Start() {
	s.lock.Lock()
	srv := s.srv
	lis := s.lis
	s.lock.Unlock()
	if srv == nil {
		return
	}
	s.logger.Info("grpc server started", "url", s.addr)
	go func() {
		if err := srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			s.notify <- fmt.Errorf("listen and serve: %w", err)
		}
	}()
}

func prepareUnixSocket(addr string) error {
	switch _, err := os.Stat(addr); {
	case errors.Is(err, fs.ErrNotExist):
		break
	case err != nil:
		return fmt.Errorf("get stat: %w", err)
	default:
		if err := os.Remove(addr); err != nil {
			return fmt.Errorf("remove file: %w", err)
		}
	}
	if err := os.MkdirAll(path.Dir(addr), os.ModePerm); err != nil {
		return fmt.Errorf("make dir: %w", err)
	}
	return nil
}

func New(
	logger logr.Logger,
	tracerProvider trace.TracerProvider,
	propagator propagation.TraceContext,
	cfg config.API,
	deps HandlerDependencyProvider,
) (*server, error) {
	network, addr := net_addr.Resolve(cfg.Address)
	if network == net_addr.UnixNetwork {
		if err := prepareUnixSocket(addr); err != nil {
			return nil, err
		}
	}
	lis, err := net.Listen(network, addr)
	if err != nil {
		return nil, fmt.Errorf("listen address: %w", err)
	}

	otelgrpc.NewServerHandler(
		otelgrpc.WithTracerProvider(tracerProvider),
		otelgrpc.WithPropagators(propagator),
	)

	srv := grpc.NewServer(
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc.ChainUnaryInterceptor(
			contextLoggerInterceptor.New(logger),
			loggerInterceptor.New(logger),
		),
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(
				otelgrpc.WithTracerProvider(tracerProvider),
				otelgrpc.WithPropagators(propagator),
			),
		),
	)

	proto.RegisterCMSServer(srv, newRouter(deps))

	return &server{
		srv:    srv,
		lis:    lis,
		logger: logger,
		addr:   cfg.Address,
		lock:   sync.Mutex{},
		notify: make(chan error, 1),
	}, nil
}
