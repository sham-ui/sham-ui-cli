package main

import (
	"net/http"
	"{{ shortName }}/test_helpers"
	"{{ shortName }}/test_helpers/asserts"
	"testing"
)

func TestSignupInvalidCSRF(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()

	resp := env.API.Request("POST", "/api/members", nil)
	asserts.Equals(t, http.StatusForbidden, resp.Response.Code, "code")
	asserts.Equals(t, "Forbidden - CSRF token not found in request\n", resp.Text(), "body")
}

func TestSignupSuccess(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	resp := env.API.Request("POST", "/api/members", map[string]interface{}{
		"Name":      "test",
		"Email":     "email",
		"Password":  "password",
		"Password2": "password",
	})
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Member created"}, resp.JSON(), "body")
}

func TestSignupInvalidData(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	resp := env.API.Request("POST", "/api/members", []string{})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Bad Request"}, resp.JSON(), "body")
}

func TestSignupPasswordMustMatch(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	resp := env.API.Request("POST", "/api/members", map[string]interface{}{
		"Name":      "test",
		"Email":     "email",
		"Password":  "1password",
		"Password2": "2password",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Bad Request", "Messages": []interface{}{"Passwords do not match."}}, resp.JSON(), "body")
}

func TestSignupEmailUnique(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	data := map[string]interface{}{
		"Name":      "test",
		"Email":     "email",
		"Password":  "password",
		"Password2": "password",
	}
	resp := env.API.Request("POST", "/api/members", data)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Member created"}, resp.JSON(), "body")

	resp = env.API.Request("POST", "/api/members", data)
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Bad Request", "Messages": []interface{}{"Email is already in use."}}, resp.JSON(), "body")
}
