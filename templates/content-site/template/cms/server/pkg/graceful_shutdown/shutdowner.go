package graceful_shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-logr/logr"
)

const shutdownTimeout = 30 * time.Second

type Shutdowner struct {
	logger logr.Logger
	tasks  []Task
	exit   chan os.Signal
	notify chan error
}

func (sh *Shutdowner) RegistryTask(task Task) *Shutdowner {
	sh.tasks = append(sh.tasks, task)
	return sh
}

func (sh *Shutdowner) Shutdown(ctx context.Context) error {
	done := make(chan bool, 1)
	go func() {
		sh.shutdownTasks(ctx)
		done <- true
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func (sh *Shutdowner) shutdownTasks(ctx context.Context) {
	for i := len(sh.tasks) - 1; i >= 0; i-- {
		task := sh.tasks[i]
		if err := task.GracefulShutdown(ctx); err != nil {
			sh.logger.Error(err, "can't graceful shutdown task", "task", task.String())
		}
	}
}

func (sh *Shutdowner) RegistryNotifier(n Notifier) *Shutdowner {
	go func() {
		if err, ok := <-n.Notify(); ok && err != nil {
			sh.logger.Error(err, "critical error, app will be shutdown", "notifier", n.String())
			sh.notify <- err
		}
	}()
	return sh
}

func (sh *Shutdowner) FailNotify(err error, msg string) {
	sh.notify <- err
	sh.logger.Error(err, msg)
}

func (sh *Shutdowner) Wait() {
	sh.logger.Info("wait signal for interrupt")
	select {
	case <-sh.exit:
		sh.logger.Info("received signal, app will be shutdown")
	case <-sh.notify:
		sh.logger.Info("received notification, app will be shutdown")
	}
	sh.logger.Info("start graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := sh.Shutdown(ctx); err != nil {
		sh.logger.Error(err, "can't graceful shutdown app")
	}
	sh.logger.Info("finish graceful shutdown")
}

func New(logger logr.Logger) *Shutdowner {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return &Shutdowner{
		logger: logger,
		tasks:  nil,
		exit:   exit,
		notify: make(chan error, 1),
	}
}
