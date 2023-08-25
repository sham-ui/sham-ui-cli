package main

import (
	"cms/test_helpers"
	"cms/test_helpers/client"
	"net/http"
	"testing"
)

func TestTag(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	_, err := env.DB.DB.Exec("DELETE FROM article_tag")
	if nil != err {
		t.Fatalf("delete from article_tag: %s", err)
	}
	_, err = env.DB.DB.Exec("DELETE FROM tag")
	if nil != err {
		t.Fatalf("delete from tag: %s", err)
	}
	_, err = env.DB.DB.Exec("INSERT INTO tag(id, name, slug) VALUES ($1, $2,$3)", "1", "test", "test")
	if nil != err {
		t.Fatalf("can't insert tag: %s", err)
	}

	env.API.ExecuteTestCases(t, []client.TestCase{
		{
			Method:                     http.MethodGet,
			URL:                        "/api/tags",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON:       `{"tags": [{"id": "1", "name": "test", "slug": "test"}]}`,
		},
		{
			Message: "empty name",
			Method:  http.MethodPost,
			URL:     "/api/tags",
			Data: map[string]interface{}{
				"name": "",
			},
			ExpectedResponseStatusCode: http.StatusBadRequest,
			ExpectedResponseJSON:       `{"Status":"Bad Request", "Messages": ["Name must not be empty."]}`,
		},
		{
			Message: "not uniq name",
			Method:  http.MethodPost,
			URL:     "/api/tags",
			Data: map[string]interface{}{
				"name": "test",
			},
			ExpectedResponseStatusCode: http.StatusBadRequest,
			ExpectedResponseJSON:       `{"Status":"Bad Request", "Messages": ["Name is already in use."]}`,
		},
		{
			Message: "create new",
			Method:  http.MethodPost,
			URL:     "/api/tags",
			Data: map[string]interface{}{
				"name": "test new",
			},
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON:       `{"Status": "Tag created"}`,
		},
		{
			Message:                    "list",
			Method:                     http.MethodGet,
			URL:                        "/api/tags",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON: `{
				"tags": [{
					"name": "test",
                    "slug": "test"
				}, {
               		"name": "test new",
                    "slug": "test-new"
				}]
			}`,
			IgnoreKeys: []string{"tags.$.id"},
		},
		{
			Message: "edit not uniq",
			Method:  http.MethodPut,
			URL:     "/api/tags/1",
			Data: map[string]interface{}{
				"name": "test",
			},
			ExpectedResponseStatusCode: http.StatusBadRequest,
			ExpectedResponseJSON:       `{"Status":"Bad Request", "Messages": ["Name is already in use."]}`,
		},
		{
			Message: "edit empty",
			Method:  http.MethodPut,
			URL:     "/api/tags/1",
			Data: map[string]interface{}{
				"name": "",
			},
			ExpectedResponseStatusCode: http.StatusBadRequest,
			ExpectedResponseJSON:       `{"Status":"Bad Request", "Messages": ["Name must not be empty."]}`,
		},
		{
			Message: "edit",
			Method:  http.MethodPut,
			URL:     "/api/tags/1",
			Data: map[string]interface{}{
				"name": "test foo",
			},
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON:       `{"Status":"Tag updated"}`,
		},
		{
			Message:                    "list",
			Method:                     http.MethodGet,
			URL:                        "/api/tags",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON: `{
				"tags": [{
					"name": "test foo",
                    "slug": "test-foo"
				}, {
               		"name": "test new",
                    "slug": "test-new"
				}]
			}`,
			IgnoreKeys: []string{"tags.$.id"},
		},
		{
			Message:                    "delete success",
			Method:                     http.MethodDelete,
			URL:                        "/api/tags/1",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON:       `{"Status":"Tag deleted"}`,
			IgnoreKeys:                 []string{"tags.$.id"},
		},
		{
			Message:                    "list",
			Method:                     http.MethodGet,
			URL:                        "/api/tags",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON: `{
				"tags": [{
               		"name": "test new",
                    "slug": "test-new"
				}]
			}`,
			IgnoreKeys: []string{"tags.$.id"},
		},
	})
}
