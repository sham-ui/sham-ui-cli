package graceful_shutdown

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"{{ shortName }}/pkg/asserts"
	"{{ shortName }}/pkg/logger"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/mock"
)

func TestShutdowner(t *testing.T) {
	// Arrange
	sh := New(testr.New(t))

	// Act
	mediumTask := NewMockTask(t)
	mediumShutdown := mediumTask.On("GracefulShutdown", mock.Anything).Return(nil).Once()
	highTask := NewMockTask(t)
	highShutdown := highTask.On("GracefulShutdown", mock.Anything).Return(nil).Once()
	highShutdown.NotBefore(mediumShutdown)

	sh.RegistryTask(highTask)
	sh.RegistryTask(mediumTask)
	err := sh.Shutdown(context.Background())

	// Assert
	asserts.NoError(t, err)
	asserts.Equals(t, 1, len(mediumTask.Calls), "medium task was shutdowned")
	asserts.Equals(t, 1, len(highTask.Calls), "medium task was shutdowned")
}

func TestShutdownerWithContext(t *testing.T) {
	// Arrange
	sh := New(testr.New(t))

	// Act
	task := NewMockTask(t)
	task.On("String").Return("task").Maybe()
	task.On("GracefulShutdown", mock.Anything).Return(nil).Maybe()
	sh.RegistryTask(task)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := sh.Shutdown(ctx)

	// Assert
	asserts.ErrorsEqual(t, errors.New("context canceled"), err)
}

func TestNotifier(t *testing.T) {
	// Arrange
	sh := New(testr.New(t))

	// Act
	task := NewMockTask(t)
	task.On("String").Return("task").Maybe()
	task.On("GracefulShutdown", mock.Anything).Return(nil).Once()
	sh.RegistryTask(task)

	notifier := NewMockNotifier(t)
	ch := make(chan error, 1)
	notifier.On("String").Return("notifier").Once()
	notifier.On("Notify").Return((<-chan error)(ch)).Once()
	sh.RegistryNotifier(notifier)

	go func() {
		ch <- errors.New("notify error")
	}()

	sh.Wait()
}

func TestNotifierFail(t *testing.T) {
	// Arrange
	sh := New(testr.New(t))

	// Act
	task := NewMockTask(t)
	task.On("GracefulShutdown", mock.Anything).Return(nil).Once()
	sh.RegistryTask(task)

	sh.FailNotify(errors.New("notify error"), "test notify error")

	sh.Wait()
}

func TestHelperProcess(t *testing.T) {
	ret := os.Getenv("GO_WANT_HELPER_PROCESS")
	if ret == "" {
		return
	}

	sh := New(logger.NewLogger(0).WithName("shutdowner"))
	{
		task := NewMockTask(t)
		var err error
		if ret == "fail" {
			err = errors.New("fail")
		}
		task.On("String").Return("task").Maybe()
		task.On("GracefulShutdown", mock.Anything).Return(err).Once()
		sh.RegistryTask(task)
	}
	{
		n := NewMockNotifier(t)
		ch := make(chan error, 1)
		n.On("String").Return("notifier").Maybe()
		n.On("Notify").Return((<-chan error)(ch)).Once()
		sh.RegistryNotifier(n)

		if ret == "notify" {
			go func() {
				ch <- errors.New("notify error")
			}()
		}
	}

	sh.Wait()
}

func TestWaitInterruptSuccess(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	//nolint:gosec
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess")
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=success"}
	cmd.Stdout = &buf
	err := cmd.Start()
	asserts.NoError(t, err)
	time.Sleep(500 * time.Millisecond)

	// Action
	err = cmd.Process.Signal(os.Interrupt)
	asserts.NoError(t, err)
	err = cmd.Wait()

	// Assert
	asserts.NoError(t, err)
	output, err := io.ReadAll(&buf)
	asserts.NoError(t, err)
	asserts.Equals(t, true, strings.Contains(string(output), `"message":"finish graceful shutdown"`), "output")
}

func TestWaitInterruptFail(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	//nolint:gosec
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess")
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=fail"}
	cmd.Stdout = &buf
	err := cmd.Start()
	asserts.NoError(t, err)
	time.Sleep(500 * time.Millisecond)

	// Action
	err = cmd.Process.Signal(os.Interrupt)
	asserts.NoError(t, err)
	err = cmd.Wait()

	// Assert
	asserts.NoError(t, err)
	output, err := io.ReadAll(&buf)
	asserts.NoError(t, err)
	asserts.Equals(t, true, strings.Contains(string(output), `"message":"received signal, app will be shutdown"`), "shutdown because receive signal")
	asserts.Equals(t, true, strings.Contains(string(output), `"message":"can't graceful shutdown task"`), "output error")
	asserts.Equals(t, true, strings.Contains(string(output), `"error":"fail"`), "output error")
	asserts.Equals(t, true, strings.Contains(string(output), `"message":"finish graceful shutdown"`), "output")
}

func TestWaitNotifier(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	//nolint:gosec
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess")
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=notify"}
	cmd.Stdout = &buf
	err := cmd.Start()
	asserts.NoError(t, err)
	time.Sleep(500 * time.Millisecond)

	// Action
	err = cmd.Wait()

	// Assert
	asserts.NoError(t, err)
	output, err := io.ReadAll(&buf)
	asserts.NoError(t, err)
	asserts.Equals(t, true, strings.Contains(string(output), `"message":"received notification, app will be shutdown"`), "shutdown because receive notify")
	asserts.Equals(t, true, strings.Contains(string(output), `"message":"critical error, app will be shutdown"`), "notify error")
	asserts.Equals(t, true, strings.Contains(string(output), `"error":"notify error"`), "output error")
	asserts.Equals(t, true, strings.Contains(string(output), `"message":"finish graceful shutdown"`), "output")
}
