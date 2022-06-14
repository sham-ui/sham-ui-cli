package main

import (
	"cms/test_helpers"
	"cms/test_helpers/asserts"
	"net/http"
	"testing"
)

func TestGetServerInfo(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateSuperUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("GET", "/api/admin/server-info", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Assert(t, len(resp.Text()) > 0, "has keys")
}

func TestGetServerInfoNonAuthorized(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	resp := env.API.Request("GET", "/api/admin/server-info", nil)
	asserts.Equals(t, http.StatusUnauthorized, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(
		t,
		[]string{},
		`{
			"Status":   "Unauthorized",
			"Messages": ["Session Expired. Log out and log back in."]
		}`,
		resp.Text(),
		"body",
	)
}

func TestGetServerInfoForNonSuperuser(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("GET", "/api/admin/server-info", nil)
	asserts.Equals(t, http.StatusForbidden, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t,
		[]string{},
		`{
			"Messages": ["Allowed only for superuser"],
			"Status":   "Forbidden"
		}`,
		resp.Text(),
		"body",
	)
}
