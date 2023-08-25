package main

import (
	"cms/config"
	"cms/test_helpers"
	"cms/test_helpers/asserts"
	"cms/test_helpers/client"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path"
	"sort"
	"testing"
	"time"
)

type article struct {
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	CategoryID  int       `json:"category_id"`
	ShortBody   string    `json:"short_body"`
	Body        string    `json:"body"`
	PublishedAt time.Time `json:"published_at"`
}

func prepareArticles(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DELETE FROM article_tag")
	if nil != err {
		t.Fatalf("delete from article_tag: %s", err)
	}
	_, err = db.Exec("DELETE FROM tag")
	if nil != err {
		t.Fatalf("delete from tag: %s", err)
	}
	_, err = db.Exec("DELETE FROM article")
	if nil != err {
		t.Fatalf("delete from article: %s", err)
	}
	_, err = db.Exec("DELETE FROM category")
	if nil != err {
		t.Fatalf("delete from category: %s", err)
	}
	_, err = db.Exec("INSERT into category(id, name, slug) VALUES(1, 'TEST', 'test')")
	if nil != err {
		t.Fatalf("insert category: %s", err)
	}
	_, err = db.Exec("INSERT into tag(id, name, slug) VALUES(1, 'first', 'first')")
	if nil != err {
		t.Fatalf("insert first tag: %s", err)
	}
	_, err = db.Exec("INSERT into tag(id, name, slug) VALUES(2, 'second', 'second')")
	if nil != err {
		t.Fatalf("insert second tag: %s", err)
	}
}

func prepareUploadDirs(t *testing.T) {
	if _, err := os.Stat(config.Upload.Path); !os.IsNotExist(err) {
		err := os.RemoveAll(config.Upload.Path)
		if nil != err {
			t.Fatalf("can't remove upload dir: %s", err)
		}
	}
	err := os.MkdirAll(config.Upload.Path, os.ModePerm)
	if nil != err {
		t.Fatalf("can't create upload dir: %s", err)
	}
}

func TestArticleList(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	prepareArticles(t, env.DB.DB)
	_, err := env.DB.DB.Exec("INSERT INTO "+
		"article(id, title, slug, category_id, short_body, body, published_at) "+
		"VALUES($1, $2, $3, $4, $5, $6, $7)",
		1,
		"test slug",
		"test-slug",
		1,
		"short body",
		"body text",
		time.Date(2022, time.March, 16, 18, 45, 0, 0, time.UTC),
	)
	if nil != err {
		t.Fatalf("insert article: %s", err)
	}
	_, err = env.DB.DB.Exec("INSERT INTO article_tag(article_id, tag_id) VALUES(1, 1)")
	if nil != err {
		t.Fatalf("insert firtst article_tag: %s", err)
	}
	_, err = env.DB.DB.Exec("INSERT INTO article_tag(article_id, tag_id) VALUES(1, 2)")
	if nil != err {
		t.Fatalf("insert second article_tag: %s", err)
	}

	env.API.ExecuteTestCases(t, []client.TestCase{
		{
			Message:                    "empty limit & offset",
			Method:                     http.MethodGet,
			URL:                        "/api/articles",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON: `{
				"articles":[
					{
						"id":"1",
						"title":"test slug",
						"slug":"test-slug",
						"category_id":"1",
						"published_at":"2022-03-17T01:45:00+07:00"
					}
				],
				"meta":{"limit":20,"offset":0,"total":1}}`,
		},
		{
			Message:                    "empty limit",
			Method:                     http.MethodGet,
			URL:                        "/api/articles?offset=0",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON: `{
				"articles":[
					{
						"id":"1",
						"title":"test slug",
						"slug":"test-slug",
						"category_id":"1",
						"published_at":"2022-03-17T01:45:00+07:00"
					}
				],
				"meta":{"limit":20,"offset":0,"total":1}}`,
		},
		{
			Message:                    "offset < 0",
			Method:                     http.MethodGet,
			URL:                        "/api/articles?offset=-1",
			ExpectedResponseStatusCode: http.StatusBadRequest,
			ExpectedResponseJSON:       `{"Messages": ["offset must be >= 0"], "Status": "Bad Request"}`,
		},
		{
			Message:                    "?limit=10&offset=0",
			Method:                     http.MethodGet,
			URL:                        "/api/articles?limit=10&offset=0",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON: `{
				"articles":[
					{
						"id":"1",
						"title":"test slug",
						"slug":"test-slug",
						"category_id":"1",
						"published_at":"2022-03-17T01:45:00+07:00"
					}
				],
				"meta":{"limit":10,"offset":0,"total":1}}`,
		},
		{
			Message:                    "?limit=10&offset=1",
			Method:                     http.MethodGet,
			URL:                        "/api/articles?limit=10&offset=1",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON: `{
				"articles": null,
				"meta":{"limit":10,"offset":1,"total":1}}`,
		},
	})
}

func TestCreateArticle(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	testCases := []struct {
		Name                 string
		Data                 map[string]interface{}
		ExpectedResponseCode int
		ExpectedResponse     string
		ExpectedTagIDs       []int
		ExpectedInDBArticle  *article
	}{
		{
			Name:                 "empty title",
			Data:                 map[string]interface{}{},
			ExpectedResponseCode: http.StatusBadRequest,
			ExpectedResponse:     `{"Messages": ["Title must not be empty."], "Status": "Bad Request"}`,
		},
		{
			Name: "success",
			Data: map[string]interface{}{
				"title":        "first article",
				"category_id":  "1",
				"short_body":   "Short body text",
				"body":         "<p>Body text</p>",
				"published_at": "2022-03-08T05:49:52.643Z",
				"tags": []map[string]string{
					{"slug": "second"},
				},
			},
			ExpectedResponseCode: http.StatusOK,
			ExpectedResponse:     `{"Status":"Article created"}`,
			ExpectedTagIDs:       []int{2},
			ExpectedInDBArticle: &article{
				Title:       "first article",
				Slug:        "first-article",
				CategoryID:  1,
				ShortBody:   "Short body text",
				Body:        "<p>Body text</p>",
				PublishedAt: time.Date(2022, time.March, 8, 12, 49, 52, 643000000, time.Local),
			},
		},
	}

	for _, testCase := range testCases {
		prepareArticles(t, env.DB.DB)
		resp := env.API.Request(http.MethodPost, "/api/articles", testCase.Data)
		asserts.Equals(t, testCase.ExpectedResponseCode, resp.Response.Code, fmt.Sprintf("%s: response code", testCase.Name))
		asserts.JSONEqualsWithoutSomeKeys(t, []string{}, testCase.ExpectedResponse, resp.Text(), fmt.Sprintf("%s: response", testCase.Name))
		rows, err := env.DB.DB.Query("select article_id, tag_id from article_tag")
		asserts.Equals(t, nil, err, fmt.Sprintf("%s: err == nil", testCase.Name))
		var ids []int
		for rows.Next() {
			var articleID int
			var tagID int
			err := rows.Scan(&articleID, &tagID)
			asserts.Equals(t, nil, err, fmt.Sprintf("%s: scan err == nil", testCase.Name))
			asserts.Equals(t, true, 0 != articleID, fmt.Sprintf("%s: articleID not 0", testCase.Name))
			ids = append(ids, tagID)
		}
		rows.Close()
		sort.Ints(ids)
		asserts.Equals(t, testCase.ExpectedTagIDs, ids, fmt.Sprintf("%s: tags", testCase.Name))
		if nil != testCase.ExpectedInDBArticle {
			articleInDB := article{}
			err := env.DB.DB.QueryRow("select title, slug, category_id, short_body, body, published_at  from article").Scan(
				&articleInDB.Title,
				&articleInDB.Slug,
				&articleInDB.CategoryID,
				&articleInDB.ShortBody,
				&articleInDB.Body,
				&articleInDB.PublishedAt,
			)
			asserts.Equals(t, nil, err, fmt.Sprintf("%s: scan article: err == nil", testCase.Name))
			asserts.Equals(t, testCase.ExpectedInDBArticle, &articleInDB, fmt.Sprintf("%s: article", testCase.Name))
		}
	}
}

func TestUpdateArticle(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	testCases := []struct {
		Name                 string
		Data                 map[string]interface{}
		ExpectedResponseCode int
		ExpectedResponse     string
		ExpectedTagIDs       []int
		ExpectedInDBArticle  *article
	}{
		{
			Name:                 "empty title",
			Data:                 map[string]interface{}{},
			ExpectedResponseCode: http.StatusBadRequest,
			ExpectedResponse:     `{"Messages": ["Title must not be empty."], "Status": "Bad Request"}`,
			ExpectedTagIDs:       []int{1, 2},
		},
		{
			Name: "success",
			Data: map[string]interface{}{
				"title":        "first article",
				"category_id":  "1",
				"short_body":   "Short body text",
				"body":         "<p>Body text</p>",
				"published_at": "2022-03-08T05:49:52.643Z",
				"tags": []map[string]string{
					{"slug": "second"},
				},
			},
			ExpectedResponseCode: http.StatusOK,
			ExpectedResponse:     `{"Status":"Article updated"}`,
			ExpectedTagIDs:       []int{2},
			ExpectedInDBArticle: &article{
				Title:       "first article",
				Slug:        "first-article",
				CategoryID:  1,
				ShortBody:   "Short body text",
				Body:        "<p>Body text</p>",
				PublishedAt: time.Date(2022, time.March, 8, 12, 49, 52, 643000000, time.Local),
			},
		},
	}

	for _, testCase := range testCases {
		prepareArticles(t, env.DB.DB)
		_, err := env.DB.DB.Exec("INSERT INTO "+
			"article(id, title, slug, category_id, short_body, body, published_at) "+
			"VALUES($1, $2, $3, $4, $5, $6, $7)",
			1,
			"test slug",
			"test-slug",
			1,
			"short body",
			"body text",
			time.Date(2022, time.March, 16, 18, 45, 0, 0, time.UTC),
		)
		if nil != err {
			t.Fatalf("insert article: %s", err)
		}
		_, err = env.DB.DB.Exec("INSERT INTO article_tag(article_id, tag_id) VALUES(1, 1)")
		if nil != err {
			t.Fatalf("insert firtst article_tag: %s", err)
		}
		_, err = env.DB.DB.Exec("INSERT INTO article_tag(article_id, tag_id) VALUES(1, 2)")
		if nil != err {
			t.Fatalf("insert second article_tag: %s", err)
		}

		resp := env.API.Request(http.MethodPut, "/api/articles/1", testCase.Data)
		asserts.Equals(t, testCase.ExpectedResponseCode, resp.Response.Code, fmt.Sprintf("%s: response code", testCase.Name))
		asserts.JSONEqualsWithoutSomeKeys(t, []string{}, testCase.ExpectedResponse, resp.Text(), fmt.Sprintf("%s: response", testCase.Name))
		rows, err := env.DB.DB.Query("select article_id, tag_id from article_tag")
		asserts.Equals(t, nil, err, fmt.Sprintf("%s: err == nil", testCase.Name))
		var ids []int
		for rows.Next() {
			var articleID int
			var tagID int
			err := rows.Scan(&articleID, &tagID)
			asserts.Equals(t, nil, err, fmt.Sprintf("%s: scan err == nil", testCase.Name))
			asserts.Equals(t, true, 0 != articleID, fmt.Sprintf("%s: articleID not 0", testCase.Name))
			ids = append(ids, tagID)
		}
		rows.Close()
		sort.Ints(ids)
		asserts.Equals(t, testCase.ExpectedTagIDs, ids, fmt.Sprintf("%s: tags", testCase.Name))
		if nil != testCase.ExpectedInDBArticle {
			articleInDB := article{}
			err := env.DB.DB.QueryRow("select title, slug, category_id, short_body, body, published_at  from article").Scan(
				&articleInDB.Title,
				&articleInDB.Slug,
				&articleInDB.CategoryID,
				&articleInDB.ShortBody,
				&articleInDB.Body,
				&articleInDB.PublishedAt,
			)
			asserts.Equals(t, nil, err, fmt.Sprintf("%s: scan article: err == nil", testCase.Name))
			asserts.Equals(t, testCase.ExpectedInDBArticle, &articleInDB, fmt.Sprintf("%s: article", testCase.Name))
		}
	}
}

func TestArticleDetail(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	prepareArticles(t, env.DB.DB)
	_, err := env.DB.DB.Exec("INSERT INTO "+
		"article(id, title, slug, category_id, short_body, body, published_at) "+
		"VALUES($1, $2, $3, $4, $5, $6, $7)",
		1,
		"test slug",
		"test-slug",
		1,
		"short body",
		"body text",
		time.Date(2022, time.March, 16, 18, 45, 0, 0, time.UTC),
	)
	if nil != err {
		t.Fatalf("insert article: %s", err)
	}
	_, err = env.DB.DB.Exec("INSERT INTO article_tag(article_id, tag_id) VALUES(1, 1)")
	if nil != err {
		t.Fatalf("insert firtst article_tag: %s", err)
	}
	_, err = env.DB.DB.Exec("INSERT INTO article_tag(article_id, tag_id) VALUES(1, 2)")
	if nil != err {
		t.Fatalf("insert second article_tag: %s", err)
	}

	env.API.ExecuteTestCases(t, []client.TestCase{
		{
			Message:                    "not found",
			Method:                     http.MethodGet,
			URL:                        "/api/articles/2",
			ExpectedResponseStatusCode: http.StatusInternalServerError,
			ExpectedResponseJSON:       `{"Status":"Internal Server Error"}`,
		},
		{
			Message:                    "success",
			Method:                     http.MethodGet,
			URL:                        "/api/articles/1",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON: `{
				"title":"test slug",
				"slug":"test-slug",
				"body": "body text",
				"short_body": "short body",
				"category_id": "1",
				"tags": ["first", "second"],
				"published_at":"2022-03-17T01:45:00+07:00"
			}`,
		},
	})
}

func TestArticleDelete(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	prepareArticles(t, env.DB.DB)
	_, err := env.DB.DB.Exec("INSERT INTO "+
		"article(id, title, slug, category_id, short_body, body, published_at) "+
		"VALUES($1, $2, $3, $4, $5, $6, $7)",
		1,
		"test slug",
		"test-slug",
		1,
		"short body",
		"body text",
		time.Date(2022, time.March, 16, 18, 45, 0, 0, time.UTC),
	)
	if nil != err {
		t.Fatalf("insert article: %s", err)
	}
	_, err = env.DB.DB.Exec("INSERT INTO article_tag(article_id, tag_id) VALUES(1, 1)")
	if nil != err {
		t.Fatalf("insert firtst article_tag: %s", err)
	}
	_, err = env.DB.DB.Exec("INSERT INTO article_tag(article_id, tag_id) VALUES(1, 2)")
	if nil != err {
		t.Fatalf("insert second article_tag: %s", err)
	}

	env.API.ExecuteTestCases(t, []client.TestCase{
		{
			Message:                    "delete",
			Method:                     http.MethodDelete,
			URL:                        "/api/articles/1",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON:       `{"Status":"Article deleted"}`,
		},
	})

	env.API.ExecuteTestCases(t, []client.TestCase{
		{
			Message:                    "get list",
			Method:                     http.MethodGet,
			URL:                        "/api/articles",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON:       `{"articles": null, "meta": {"limit": 20, "offset": 0, "total": 0}}`,
		},
	})
}

func TestArticleUploadImages(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	prepareUploadDirs(t)

	gifFile, err := os.ReadFile(path.Join("testdata", "Sorting_quicksort_anim.gif"))
	if nil != err {
		t.Fatalf("can't read test gif file: %s", err)
	}

	testCases := []struct {
		Message                    string
		Request                    *client.MulipartRequest
		ExpectedResponseStatusCode int
		ExpectedResponse           string
	}{
		{
			Message:                    "empty request",
			Request:                    &client.MulipartRequest{},
			ExpectedResponseStatusCode: http.StatusBadRequest,
			ExpectedResponse:           `{"Status": "Bad Request"}`,
		},
		{
			Message: "empty content",
			Request: &client.MulipartRequest{
				FileField: "file-0",
				FilePath:  "image.jpg",
			},
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponse: `{
				"errorMessage":"",
				"result":[{
					"url":"/assets/d41d8cd98f00b204e9800998ecf8427e.jpg",
					"name":"d41d8cd98f00b204e9800998ecf8427e.jpg",
					"size":"0"
				}]
			}`,
		},
		{
			Message: "success",
			Request: &client.MulipartRequest{
				FileField: "file-0",
				FilePath:  "image.png",
				Content:   gifFile,
			},
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponse: `{
				"errorMessage":"",
				"result":[{
					"url":"/assets/d4e5d0a778dba725091d8317e6bac939.png",
					"name":"d4e5d0a778dba725091d8317e6bac939.png",
					"size":"93016"
				}]
			}`,
		},
	}

	for _, test := range testCases {
		resp := env.API.RequestMultiPart(http.MethodPost, "/api/upload-image", test.Request)
		asserts.Equals(t, test.ExpectedResponseStatusCode, resp.Response.Code, fmt.Sprintf("%s: http code", test.Message))
		asserts.JSONEqualsWithoutSomeKeys(t, []string{}, test.ExpectedResponse, resp.Text(), fmt.Sprintf("%s: body", test.Message))
	}
}

func TestArticleFetchImages(t *testing.T) {
	env := test_helpers.NewTestEnv(t)
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	prepareUploadDirs(t)

	gifFile, err := os.ReadFile(path.Join("testdata", "Sorting_quicksort_anim.gif"))
	if nil != err {
		t.Fatalf("can't read test gif file: %s", err)
	}
	f, err := os.Create(path.Join(config.Upload.Path, "d4e5d0a778dba725091d8317e6bac939.gif"))
	if nil != err {
		t.Fatalf("can't create file: %s", err)
	}
	_, err = f.Write(gifFile)
	if nil != err {
		t.Fatalf("can't write file: %s", err)
	}
	f.Close()

	testCases := []struct {
		Message                    string
		Path                       string
		ExpectedResponseStatusCode int
		ExpectedResponse           []byte
	}{
		{
			Message:                    "not found",
			Path:                       "/assets/not-found.jpg",
			ExpectedResponseStatusCode: http.StatusNotFound,
			ExpectedResponse:           []byte("404 page not found\n"),
		},
		{
			Message:                    "success",
			Path:                       "/assets/d4e5d0a778dba725091d8317e6bac939.gif",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponse:           gifFile,
		},
	}

	for _, test := range testCases {
		resp := env.API.Request(http.MethodGet, test.Path, nil)
		asserts.Equals(t, test.ExpectedResponseStatusCode, resp.Response.Code, fmt.Sprintf("%s: http code", test.Message))
		asserts.Equals(t, test.ExpectedResponse, resp.Response.Body.Bytes(), fmt.Sprintf("%s: body", test.Message))
	}
}
