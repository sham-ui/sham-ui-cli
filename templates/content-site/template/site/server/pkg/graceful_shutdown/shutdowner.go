package graceful_shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-logr/logr"
)

const shutdownTimeout = 30 * time.Second

type taskWithPriority struct {
	task     Task
	priority Priority
}

type Shutdowner struct {
	logger logr.Logger
	tasks  []taskWithPriority
	exit   chan os.Signal
	notify chan error
}

func (sh *Shutdowner) RegistryTask(task Task) *Shutdowner {
	sh.registryWithPriority(task, PriorityMedium)
	return sh
}

func (sh *Shutdowner) RegistryHighTask(task Task) *Shutdowner {
	sh.registryWithPriority(task, PriorityHigh)
	return sh
}

func (sh *Shutdowner) registryWithPriority(task Task, priority Priority) {
	sh.tasks = append(sh.tasks, taskWithPriority{
		task:     task,
		priority: priority,
	})
}

func (sh *Shutdowner) Shutdown(ctx context.Context) error {
	done := make(chan bool, 1)
	go func() {
		sh.shutdownWithPriority(ctx, PriorityMedium)
		sh.shutdownWithPriority(ctx, PriorityHigh)
		done <- true
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func (sh *Shutdowner) shutdownWithPriority(ctx context.Context, priority Priority) {
	stopped := &sync.WaitGroup{}
	stopped.Add(len(sh.tasks))
	for _, task := range sh.tasks {
		if task.priority == priority {
			go func(task Task) {
				if err := task.GracefulShutdown(ctx); err != nil {
					sh.logger.Error(err, "can't graceful shutdown task", "task", task.String())
				}
				stopped.Done()
			}(task.task)
		} else {
			stopped.Done()
		}
	}
	stopped.Wait()
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
