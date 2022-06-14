package main

import (
	"cms/test_helpers"
	"cms/test_helpers/asserts"
	"net/http"
	"testing"
)

func TestCsrfToken(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()

	resp := env.API.Request("GET", "/api/csrftoken", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Equals(t, "", resp.Text(), "text")
	asserts.Equals(t, 88, len(resp.Response.Header().Get("X-CSRF-Token")), "token")
}

func TestSessionNotExists(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	resp := env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusUnauthorized, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Status": "Unauthorized", "Messages": ["Session Expired. Log out and log back in."]}`, resp.Text(), "body")
}

func TestSessionExists(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Name": "test", "Email": "email", "IsSuperuser": false}`, resp.Text(), "body")
}

func TestSuperuserSession(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateSuperUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Name": "test", "Email": "email", "IsSuperuser": true}`, resp.Text(), "body")
}
