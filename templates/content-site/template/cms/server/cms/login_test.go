package main

import (
	"cms/test_helpers"
	"cms/test_helpers/asserts"
	"net/http"
	"testing"
)

func TestLoginInvalidCSRF(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()

	resp := env.API.Request("POST", "/api/login", nil)
	asserts.Equals(t, http.StatusForbidden, resp.Response.Code, "code")
	asserts.Equals(t, "Forbidden - CSRF token not found in request\n", resp.Text(), "body")
}

func TestLoginSuccess(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()

	resp := env.API.Request("POST", "/api/login", map[string]interface{}{
		"Email":    "email",
		"Password": "password",
	})
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t,
		[]string{},
		`{"Status": "OK", "IsSuperuser": false, "Name": "test", "Email": "email"}`,
		resp.Text(),
		"body",
	)
}

func TestLoginIncorrectPassword(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()

	resp := env.API.Request("POST", "/api/login", map[string]interface{}{
		"Email":    "email",
		"Password": "incorrectPassword",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(
		t,
		[]string{},
		`{"Status": "Bad Request", "Messages": ["Incorrect username or password"]}`,
		resp.Text(),
		"body",
	)
}

func TestLoginIncorrectEmail(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()

	resp := env.API.Request("POST", "/api/login", map[string]interface{}{
		"Email":    "incorrectemail",
		"Password": "password",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(
		t,
		[]string{},
		`{"Status": "Bad Request"}`,
		resp.Text(),
		"body",
	)
}
