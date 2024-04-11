package session

import (
	"cms/config"
	"cms/internal/model"
	"cms/pkg/asserts"
	"cms/pkg/graceful_shutdown"
	"cms/pkg/postgres/testingdb"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
)

// Compile-time check that service implements graceful_shutdown.Task.
var _ graceful_shutdown.Task = (*service)(nil)

var member = model.Member{
	ID:          "42",
	Email:       "test@example.com",
	Name:        "tester",
	IsSuperuser: true,
}

func newService(t *testing.T) *service {
	t.Helper()
	log := testr.New(t)
	cfg, err := config.LoadConfiguration(log, "../../../testdata/config.cfg")
	asserts.NoError(t, err)
	srv, err := New(testingdb.Connect(t), cfg.Session.Secret, cfg.Session.Domain, 5*time.Minute)
	asserts.NoError(t, err)
	t.Cleanup(func() {
		asserts.NoError(t, srv.GracefulShutdown(context.Background()))
	})
	return srv
}

func newRequestWithCookie(resp *httptest.ResponseRecorder) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Cookie", resp.Result().Header.Get("Set-Cookie")) //nolint:bodyclose
	return req
}

func TestService_Create(t *testing.T) {
	// Arrange
	srv := newService(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	// Action
	err := srv.Create(resp, req, &member)

	// Assert
	asserts.NoError(t, err)
	asserts.Equals(t, true, len(resp.Result().Header.Get("Set-Cookie")) > 0) //nolint:bodyclose
}

func TestService_Get(t *testing.T) {
	// Arrange
	srv := newService(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	err := srv.Create(resp, req, &member)
	asserts.NoError(t, err)

	// Action
	req = newRequestWithCookie(resp)
	sess, err := srv.Get(req)

	// Assert
	asserts.NoError(t, err)
	asserts.Equals(t, &model.Session{
		MemberID:    member.ID,
		Email:       member.Email,
		Name:        member.Name,
		IsSuperuser: member.IsSuperuser,
	}, sess)
}

func TestService_Delete(t *testing.T) {
	// Arrange
	srv := newService(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	err := srv.Create(resp, req, &member)
	asserts.NoError(t, err)

	// Action
	req = newRequestWithCookie(resp)
	resp = httptest.NewRecorder()
	err = srv.Delete(resp, req)
	asserts.NoError(t, err)
	_, err = srv.Get(req)

	// Assert
	asserts.ErrorsEqual(t, errors.New("session not authenticated"), err)
}

func TestService_UpdateEmail(t *testing.T) {
	// Arrange
	srv := newService(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	err := srv.Create(resp, req, &member)
	asserts.NoError(t, err)

	// Action
	req = newRequestWithCookie(resp)
	resp = httptest.NewRecorder()
	err = srv.UpdateEmail(resp, req, "changed@example.com")

	// Assert
	asserts.NoError(t, err)
	req = newRequestWithCookie(resp)
	sess, err := srv.Get(req)
	asserts.NoError(t, err)
	asserts.Equals(t, &model.Session{
		MemberID:    member.ID,
		Email:       "changed@example.com",
		Name:        member.Name,
		IsSuperuser: member.IsSuperuser,
	}, sess, "email changed")
}

func TestService_UpdateName(t *testing.T) {
	// Arrange
	srv := newService(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	err := srv.Create(resp, req, &member)
	asserts.NoError(t, err)

	// Action
	req = newRequestWithCookie(resp)
	resp = httptest.NewRecorder()
	err = srv.UpdateName(resp, req, "changed-tester")

	// Assert
	asserts.NoError(t, err)
	req = newRequestWithCookie(resp)
	sess, err := srv.Get(req)
	asserts.NoError(t, err)
	asserts.Equals(t, &model.Session{
		MemberID:    member.ID,
		Email:       member.Email,
		Name:        "changed-tester",
		IsSuperuser: member.IsSuperuser,
	}, sess, "name changed")
}
