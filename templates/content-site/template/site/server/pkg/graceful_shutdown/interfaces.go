package graceful_shutdown

import "context"

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name Task --inpackage --testonly
type Task interface {
	String() string
	GracefulShutdown(ctx context.Context) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name Notifier --inpackage --testonly
type Notifier interface {
	String() string
	Notify() <-chan error
}
