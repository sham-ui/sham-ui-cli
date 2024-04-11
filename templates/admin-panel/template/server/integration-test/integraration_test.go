package integration

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os/exec"
	"{{ shortName }}/config"
	"{{ shortName }}/internal/model"
	"{{ shortName }}/internal/repository/postgres/member"
	"{{ shortName }}/internal/service/password"
	"{{ shortName }}/pkg/asserts"
	"{{ shortName }}/pkg/logger"
	"{{ shortName }}/pkg/postgres"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
)

func runApp(t *testing.T) *exec.Cmd {
	t.Helper()
	var buf bytes.Buffer
	cmd := exec.Command("./app")
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	asserts.NoError(t, cmd.Start())
	t.Cleanup(func() {
		_ = cmd.Process.Signal(syscall.SIGTERM)
		_ = cmd.Wait()
	})

	timeout := time.After(5 * time.Second)
	for {
		select {
		case <-timeout:
			t.Fatalf("timeout, expected string not found: %s", buf.String())
		default:
			time.Sleep(10 * time.Millisecond)
			if strings.Contains(buf.String(), "wait signal for interrupt") {
				return cmd
			}
		}
	}
}

func createSuperuser(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	log := logger.NewLogger(128)
	cfg, err := config.LoadConfiguration(log, "config.cfg")
	asserts.NoError(t, err)
	db, err := postgres.New(log, cfg.Database.URL())
	asserts.NoError(t, err)
	pass, err := password.New().Hash(ctx, "password")
	asserts.NoError(t, err)
	repo := member.NewRepository(db)
	err = repo.Create(ctx, model.MemberWithPassword{
		Member: model.Member{ //nolint:exhaustruct
			Email:       "root",
			Name:        "Superuser",
			IsSuperuser: true,
		},
		HashedPassword: pass,
	})
	if !errors.Is(err, model.ErrMemberEmailAlreadyExists) {
		asserts.NoError(t, err)
	}
	mem, err := repo.GetByEmail(ctx, "root")
	asserts.NoError(t, err)
	t.Cleanup(func() {
		asserts.NoError(t, repo.Delete(ctx, mem.ID))
	})
}

func TestSmoke(t *testing.T) {
	// Action
	cmd := runApp(t)
	stdout := (cmd.Stdout).(*bytes.Buffer).String()

	// Assert
	asserts.Contains(t, stdout, "http server started")
}

func TestGracefulShutdown(t *testing.T) {
	// Arrange
	cmd := runApp(t)

	// Action
	err := cmd.Process.Signal(syscall.SIGTERM)
	asserts.NoError(t, err)
	err = cmd.Wait()
	asserts.NoError(t, err)
	stdout := (cmd.Stdout).(*bytes.Buffer).String()

	// Assert
	asserts.NotContains(t, stdout, "can't graceful shutdown app")
	asserts.Contains(t, stdout, "finish graceful shutdown")
}

func Test_HTTP(t *testing.T) {
	cfg, err := config.LoadConfiguration(testr.New(t).V(1), "config.cfg")
	asserts.NoError(t, err)
	cmd := runApp(t)
	stdout := cmd.Stdout.(*bytes.Buffer)

	testCases := []struct {
		Name         string
		URL          string
		ExpectedCode int
		ExpectedBody string
	}{
		{
			Name:         "home",
			URL:          "/",
			ExpectedCode: http.StatusOK,
			ExpectedBody: `<script>System.import('/bundle.js');</script>`,
		},
		{
			Name:         "login",
			URL:          "/login",
			ExpectedCode: http.StatusOK,
			ExpectedBody: `<script>System.import('/bundle.js');</script>`,
		},
		{
			Name:         "login strip slash",
			URL:          "/login/",
			ExpectedCode: http.StatusOK,
			ExpectedBody: `<script>System.import('/bundle.js');</script>`,
		},
		{
			Name:         "bundle.js",
			URL:          "/bundle.js",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "sourceMappingURL=bundle.js.map",
		},
		{
			Name:         "csrf token",
			URL:          "/api/csrftoken",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			// Arrange
			stdout.Reset()

			// Action
			//nolint:bodyclose,noctx
			resp, err := (&http.Client{Timeout: 5 * time.Second}).Get(cfg.Server.URL() + test.URL)
			asserts.NoError(t, err)
			b, err := io.ReadAll(resp.Body)
			asserts.NoError(t, err)

			// Assert
			asserts.Equals(t, test.ExpectedCode, resp.StatusCode, "status code")
			asserts.Contains(t, string(b), test.ExpectedBody)
			asserts.Contains(t, stdout.String(), test.URL)
		})
	}
}

func Test_HTTPWithCSRF(t *testing.T) {
	cmd := runApp(t)
	createSuperuser(t)
	stdout := cmd.Stdout.(*bytes.Buffer)

	testCases := []struct {
		Name         string
		URL          string
		Method       string
		Body         []byte
		ExpectedCode int
		ExpectedBody string
	}{
		{
			Name:         "validate session",
			URL:          "/api/validsession",
			Method:       http.MethodGet,
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: "not authenticated",
		},
		{
			Name:   "login",
			URL:    "/api/login",
			Method: http.MethodPost,
			Body: []byte(`{
				"email": "root",
				"password": "password"
			}`),
			ExpectedCode: http.StatusOK,
			ExpectedBody: `{"Status":"OK","Name":"Superuser","Email":"root","IsSuperuser":true}`,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			// Arrange
			stdout.Reset()
			const csrfToken = "X-Csrf-Token" //nolint:gosec
			cfg, err := config.LoadConfiguration(testr.New(t).V(1), "config.cfg")
			asserts.NoError(t, err)
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Get(cfg.Server.URL() + "/api/csrftoken") //nolint:bodyclose,noctx
			asserts.NoError(t, err)
			token := resp.Header.Get(csrfToken)
			cookies := resp.Cookies()

			// Action
			//nolint:noctx
			req, err := http.NewRequest(test.Method, cfg.Server.URL()+test.URL, bytes.NewReader(test.Body))
			asserts.NoError(t, err)
			for _, cookie := range cookies {
				req.AddCookie(cookie)
			}
			req.Header.Set(csrfToken, token)
			resp, err = client.Do(req) //nolint:bodyclose
			asserts.NoError(t, err)
			b, err := io.ReadAll(resp.Body)
			asserts.NoError(t, err)

			// Assert
			asserts.Equals(t, test.ExpectedCode, resp.StatusCode, "status code")
			asserts.Contains(t, string(b), test.ExpectedBody)
			asserts.Contains(t, stdout.String(), test.URL)
		})
	}
}
