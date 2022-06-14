package main

import (
	"net/http"
	"{{ shortName }}/test_helpers"
	"{{ shortName }}/test_helpers/asserts"
	"testing"
)

func TestLoginInvalidCSRF(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()

	resp := env.API.Request("POST", "/api/login", nil)
	asserts.Equals(t, http.StatusForbidden, resp.Response.Code, "code")
	asserts.Equals(t, "Forbidden - CSRF token not found in request\n", resp.Text(), "body")
}

func TestLoginSuccess(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()

	resp := env.API.Request("POST", "/api/login", map[string]interface{}{
		"Email":    "email",
		"Password": "password",
	})
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "OK", "IsSuperuser": false, "Name": "test", "Email": "email"}, resp.JSON(), "body")
}

func TestLoginIncorrectPassword(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()

	resp := env.API.Request("POST", "/api/login", map[string]interface{}{
		"Email":    "email",
		"Password": "incorrectPassword",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Bad Request", "Messages": []interface{}{"Incorrect username or password"}}, resp.JSON(), "body")
}

func TestLoginIncorrectEmail(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()

	resp := env.API.Request("POST", "/api/login", map[string]interface{}{
		"Email":    "incorrectemail",
		"Password": "password",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Bad Request"}, resp.JSON(), "body")
}
