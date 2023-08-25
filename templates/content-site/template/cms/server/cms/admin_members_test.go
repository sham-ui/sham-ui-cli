package main

import (
	"cms/test_helpers"
	"cms/test_helpers/client"
	"net/http"
	"testing"
)

func TestMembers(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateSuperUser()
	env.API.GetCSRF()
	env.API.Login()

	env.API.ExecuteTestCases(t, []client.TestCase{
		{
			Method:                     http.MethodGet,
			URL:                        "/api/admin/members?limit=20",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON: `{
				"members": [{
					"Email": "email",
                    "ID": "1",
					"IsSuperuser": true,
					"Name": "test"
				}],
                "meta": {
					"limit": 20,
                    "offset": 0,
                    "total": 1
                }
			}`,
		},
		{
			Method: http.MethodPost,
			URL:    "/api/admin/members",
			Data: map[string]interface{}{
				"name":         "new_user",
				"email":        "new_user@email.com",
				"is_superuser": true,
				"password":     "test",
			},
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON:       `{"Status": "Member created"}`,
		},
		{
			Message: "Empty email",
			Method:  http.MethodPost,
			URL:     "/api/admin/members",
			Data: map[string]interface{}{
				"name":         "new_user",
				"email":        "",
				"is_superuser": true,
				"password":     "test",
			},
			ExpectedResponseStatusCode: http.StatusBadRequest,
			ExpectedResponseJSON: `{
				"Status":   "Bad Request",
				"Messages": ["Email must not be empty"]
			}`,
		},
		{
			Message: "Empty name",
			Method:  http.MethodPost,
			URL:     "/api/admin/members",
			Data: map[string]interface{}{
				"name":         "",
				"email":        "new_user1@email.com",
				"is_superuser": true,
				"password":     "test",
			},
			ExpectedResponseStatusCode: http.StatusBadRequest,
			ExpectedResponseJSON: `{
				"Status":   "Bad Request",
				"Messages": ["Name must not be empty."]
			}`,
		},
		{
			Message: "Not unique email",
			Method:  http.MethodPost,
			URL:     "/api/admin/members",
			Data: map[string]interface{}{
				"name":         "name",
				"email":        "email",
				"is_superuser": true,
				"password":     "test",
			},
			ExpectedResponseStatusCode: http.StatusBadRequest,
			ExpectedResponseJSON: `{
				"Status":   "Bad Request",
				"Messages": ["Email is already in use."]
			}`,
		},
		{
			Method: http.MethodPut,
			URL:    "/api/admin/members/1",
			Data: map[string]interface{}{
				"name":         "new_user2",
				"email":        "new_user2@email.com",
				"is_superuser": false,
			},
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON:       `{"Status": "Member updated"}`,
		},
		{
			Message: "Only email",
			Method:  http.MethodPut,
			URL:     "/api/admin/members/1",
			Data: map[string]interface{}{
				"name":         "new_user3",
				"email":        "new_user2@email.com",
				"is_superuser": false,
			},
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON:       `{"Status": "Member updated"}`,
		},
		{
			Message: "Empty email",
			Method:  http.MethodPut,
			URL:     "/api/admin/members/1",
			Data: map[string]interface{}{
				"name":         "new_user2",
				"email":        "",
				"is_superuser": false,
			},
			ExpectedResponseStatusCode: http.StatusBadRequest,
			ExpectedResponseJSON: `{
				"Status":   "Bad Request",
				"Messages": ["Email must not be empty"]
			}`,
		},
		{
			Message: "Empty name",
			Method:  http.MethodPut,
			URL:     "/api/admin/members/1",
			Data: map[string]interface{}{
				"name":         "",
				"email":        "new_user2@email.com",
				"is_superuser": false,
			},
			ExpectedResponseStatusCode: http.StatusBadRequest,
			ExpectedResponseJSON: `{
				"Status":   "Bad Request",
				"Messages": ["Name must not be empty."]
			}`,
		},
		{
			Method: http.MethodPut,
			URL:    "/api/admin/members/1/password",
			Data: map[string]interface{}{
				"pass1": "pass1",
				"pass2": "pass1",
			},
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON:       `{"Status": "Password updated"}`,
		},
		{
			Message: "Empty pass1",
			Method:  http.MethodPut,
			URL:     "/api/admin/members/1/password",
			Data: map[string]interface{}{
				"pass1": "",
				"pass2": "pass2",
			},
			ExpectedResponseStatusCode: http.StatusBadRequest,
			ExpectedResponseJSON: `{
				"Status": "Bad Request",
				"Messages": [
					"Password must have more than 0 characters.",
					"Passwords don't match."
				]
			}`,
		},
		{
			Message: "empty pass2",
			Method:  http.MethodPut,
			URL:     "/api/admin/members/1/password",
			Data: map[string]interface{}{
				"pass1": "pass1",
				"pass2": "",
			},
			ExpectedResponseStatusCode: http.StatusBadRequest,
			ExpectedResponseJSON: `{
				"Status": "Bad Request",
				"Messages": [
					"Password must have more than 0 characters.",
					"Passwords don't match."
				]
			}`,
		},
		{
			Message: "passwords don't match",
			Method:  http.MethodPut,
			URL:     "/api/admin/members/1/password",
			Data: map[string]interface{}{
				"pass1": "pass1",
				"pass2": "pass2",
			},
			ExpectedResponseStatusCode: http.StatusBadRequest,
			ExpectedResponseJSON: `{
				"Status": "Bad Request",
				"Messages": [
					"Passwords don't match."
				]
			}`,
		},
		{
			Message:                    "delete success",
			Method:                     http.MethodDelete,
			URL:                        "/api/admin/members/1",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON:       `{"Status": "Member deleted"}`,
		},
		{
			Message:                    "member not found",
			Method:                     http.MethodDelete,
			URL:                        "/api/admin/members/0",
			ExpectedResponseStatusCode: http.StatusBadRequest,
			ExpectedResponseJSON: `{
				"Status": "Bad Request",
				"Messages": [
					"Member not exists."
				]
			}`,
		},
	})
}

func TestUpdateMemberEmailNotUnique(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateSuperUser()
	env.API.GetCSRF()
	env.API.Login()

	env.DB.DB.Exec("INSERT INTO public.members (name, email, password, is_superuser) VALUES ('test', 'email2', '$2a$14$QMQH3E2UyfIKTFvLfguQPOmai96AncIV.1bLbcd5huTG8gZxNfAyO', TRUE)")

	env.API.ExecuteTestCases(t, []client.TestCase{

		{
			Method: http.MethodPut,
			URL:    "/api/admin/members/1",
			Data: map[string]interface{}{
				"name":         "new_user2",
				"email":        "email2",
				"is_superuser": false,
			},
			ExpectedResponseStatusCode: http.StatusBadRequest,
			ExpectedResponseJSON: `{
				"Status":   "Bad Request",
				"Messages": ["Email is already in use."]
			}`,
		},
	})
}

func TestMembersNonAuthorized(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	env.API.ExecuteTestCases(t, []client.TestCase{
		{
			Method:                     http.MethodGet,
			URL:                        "/api/admin/members",
			ExpectedResponseStatusCode: http.StatusUnauthorized,
			ExpectedResponseJSON: `{
				"Status":   "Unauthorized",
				"Messages": ["Session Expired. Log out and log back in."]
			}`,
		},
		{
			Method:                     http.MethodPost,
			URL:                        "/api/admin/members",
			ExpectedResponseStatusCode: http.StatusUnauthorized,
			ExpectedResponseJSON: `{
				"Status":   "Unauthorized",
				"Messages": ["Session Expired. Log out and log back in."]
			}`,
		},
		{
			Method:                     http.MethodPut,
			URL:                        "/api/admin/members/1",
			ExpectedResponseStatusCode: http.StatusUnauthorized,
			ExpectedResponseJSON: `{
				"Status":   "Unauthorized",
				"Messages": ["Session Expired. Log out and log back in."]
			}`,
		},
		{
			Method:                     http.MethodPut,
			URL:                        "/api/admin/members/1/password",
			ExpectedResponseStatusCode: http.StatusUnauthorized,
			ExpectedResponseJSON: `{
				"Status":   "Unauthorized",
				"Messages": ["Session Expired. Log out and log back in."]
			}`,
		},
		{
			Method:                     http.MethodDelete,
			URL:                        "/api/admin/members/1",
			ExpectedResponseStatusCode: http.StatusUnauthorized,
			ExpectedResponseJSON: `{
				"Status":   "Unauthorized",
				"Messages": ["Session Expired. Log out and log back in."]
			}`,
		},
	})
}

func TestMembersForNonSuperuser(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	env.API.ExecuteTestCases(t, []client.TestCase{
		{
			Method:                     http.MethodGet,
			URL:                        "/api/admin/members",
			ExpectedResponseStatusCode: http.StatusForbidden,
			ExpectedResponseJSON: `{
				"Messages": ["Allowed only for superuser"],
				"Status":   "Forbidden"
			}`,
		},
		{
			Method:                     http.MethodPost,
			URL:                        "/api/admin/members",
			ExpectedResponseStatusCode: http.StatusForbidden,
			ExpectedResponseJSON: `{
				"Messages": ["Allowed only for superuser"],
				"Status":   "Forbidden"
			}`,
		},
		{
			Method:                     http.MethodPut,
			URL:                        "/api/admin/members/1",
			ExpectedResponseStatusCode: http.StatusForbidden,
			ExpectedResponseJSON: `{
				"Messages": ["Allowed only for superuser"],
				"Status":   "Forbidden"
			}`,
		},
		{
			Method:                     http.MethodPut,
			URL:                        "/api/admin/members/1/password",
			ExpectedResponseStatusCode: http.StatusForbidden,
			ExpectedResponseJSON: `{
				"Messages": ["Allowed only for superuser"],
				"Status":   "Forbidden"
			}`,
		},
		{
			Method:                     http.MethodDelete,
			URL:                        "/api/admin/members/1",
			ExpectedResponseStatusCode: http.StatusForbidden,
			ExpectedResponseJSON: `{
				"Messages": ["Allowed only for superuser"],
				"Status":   "Forbidden"
			}`,
		},
	})
}
