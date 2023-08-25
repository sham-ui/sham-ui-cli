package main

import (
	"cms/test_helpers"
	"cms/test_helpers/asserts"
	"net/http"
	"testing"
)

func TestGetSuMemberListBundle(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateSuperUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("GET", "/su_members.js", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Equals(t, true, len(resp.Response.Body.Bytes()) > 0, "has response")
}

func TestGetSuMemberListBundleNonAutorized(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	resp := env.API.Request("GET", "/su_members.js", nil)
	asserts.Equals(t, http.StatusUnauthorized, resp.Response.Code, "code")
	asserts.Equals(t, "Unauthorized\n", resp.Text(), "text")
}

func TestGetSuMemberListBundleForNonSuperuser(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("GET", "/su_members.js", nil)
	asserts.Equals(t, http.StatusForbidden, resp.Response.Code, "code")
	asserts.Equals(t, "Forbidden\n", resp.Text(), "text")
}
