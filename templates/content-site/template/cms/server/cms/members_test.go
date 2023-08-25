package main

import (
	"cms/test_helpers"
	"cms/test_helpers/asserts"
	"net/http"
	"testing"
)

func TestUpdateNameSuccess(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/name", map[string]interface{}{
		"NewName": "edited test name",
	})
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Status": "Name updated"}`, resp.Text(), "body")

	resp = env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Name": "edited test name", "Email": "email", "IsSuperuser": false}`, resp.Text(), "body")
}

func TestUpdateNameUnauthtorized(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()

	resp := env.API.Request("PUT", "/api/members/name", map[string]interface{}{
		"NewName": "edited test name",
	})
	asserts.Equals(t, http.StatusUnauthorized, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Status": "Unauthorized", "Messages": ["Session Expired. Log out and log back in."]}`, resp.Text(), "body")
}

func TestUpdateNameShortName(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/name", map[string]interface{}{
		"NewName": "",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Status": "Bad Request", "Messages": ["Name must have more than 0 characters."]}`, resp.Text(), "body")

	resp = env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Name": "test", "Email": "email", "IsSuperuser": false}`, resp.Text(), "body")
}

func TestUpdateEmailSuccess(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/email", map[string]interface{}{
		"NewEmail1": "newemail@test.com",
		"NewEmail2": "newemail@test.com",
	})
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Status": "Email updated"}`, resp.Text(), "body")

	resp = env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Name": "test", "Email": "newemail@test.com", "IsSuperuser": false}`, resp.Text(), "body")
}

func TestUpdateEmailUnauthtorized(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()

	resp := env.API.Request("PUT", "/api/members/email", map[string]interface{}{
		"NewEmail1": "newemail@test.com",
		"NewEmail2": "newemail@test.com",
	})
	asserts.Equals(t, http.StatusUnauthorized, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Status": "Unauthorized", "Messages": ["Session Expired. Log out and log back in."]}`, resp.Text(), "body")
}

func TestUpdateEmailShort(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/email", map[string]interface{}{
		"NewEmail1": "",
		"NewEmail2": "newemail@test.com",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Status": "Bad Request", "Messages": ["Email must have more than 0 characters.", "Emails don't match."]}`, resp.Text(), "body")

	resp = env.API.Request("PUT", "/api/members/email", map[string]interface{}{
		"NewEmail1": "newemail@test.com",
		"NewEmail2": "",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Status": "Bad Request", "Messages": ["Email must have more than 0 characters.", "Emails don't match."]}`, resp.Text(), "body")

	resp = env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Name": "test", "Email": "email", "IsSuperuser": false}`, resp.Text(), "body")
}

func TestUpdateEmailNotMatch(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/email", map[string]interface{}{
		"NewEmail1": "email1",
		"NewEmail2": "email2",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Status": "Bad Request", "Messages": ["Emails don't match."]}`, resp.Text(), "body")

	resp = env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Name": "test", "Email": "email", "IsSuperuser": false}`, resp.Text(), "body")
}

func TestUpdateEmailNotUnique(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.DB.DB.Exec("INSERT INTO public.members (id, name, email, password) VALUES (2, 'test', 'email1', '$2a$14$QMQH3E2UyfIKTFvLfguQPOmai96AncIV.1bLbcd5huTG8gZxNfAyO')")
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/email", map[string]interface{}{
		"NewEmail1": "email1",
		"NewEmail2": "email1",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Status": "Bad Request", "Messages": ["Email is already in use."]}`, resp.Text(), "body")

	resp = env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Name": "test", "Email": "email", "IsSuperuser": false}`, resp.Text(), "body")
}

func TestUpdatePasswordSuccess(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/password", map[string]interface{}{
		"NewPassword1": "newpass",
		"NewPassword2": "newpass",
	})
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Status": "Password updated"}`, resp.Text(), "body")

	env.API.ResetCSRF()
	env.API.ResetCookies()
	env.API.GetCSRF()

	resp = env.API.Request("POST", "/api/login", map[string]interface{}{
		"Email":    "email",
		"Password": "password",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Status": "Bad Request", "Messages": ["Incorrect username or password"]}`, resp.Text(), "body")

	resp = env.API.Request("POST", "/api/login", map[string]interface{}{
		"Email":    "email",
		"Password": "newpass",
	})
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
}

func TestUpdatePasswordUnauthtorized(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	resp := env.API.Request("PUT", "/api/members/password", map[string]interface{}{
		"NewPassword1": "newpass",
		"NewPassword2": "newpass",
	})
	asserts.Equals(t, http.StatusUnauthorized, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Status": "Unauthorized", "Messages": ["Session Expired. Log out and log back in."]}`, resp.Text(), "body")
}

func TestUpdatePasswordShort(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/password", map[string]interface{}{
		"NewPassword1": "",
		"NewPassword2": "newpass",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Status": "Bad Request", "Messages": ["Password must have more than 0 characters.", "Passwords don't match."]}`, resp.Text(), "body")

	resp = env.API.Request("PUT", "/api/members/password", map[string]interface{}{
		"NewPassword1": "newpass",
		"NewPassword2": "",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Status": "Bad Request", "Messages": ["Password must have more than 0 characters.", "Passwords don't match."]}`, resp.Text(), "body")
}

func TestUpdatePasswordNotMatch(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/password", map[string]interface{}{
		"NewPassword1": "newpass1",
		"NewPassword2": "newpass2",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.JSONEqualsWithoutSomeKeys(t, nil, `{"Status": "Bad Request", "Messages": ["Passwords don't match."]}`, resp.Text(), "body")
}
