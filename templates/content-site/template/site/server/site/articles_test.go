package main

import (
	"encoding/json"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"site/proto"
	"site/test_helpers"
	"site/test_helpers/client"
	"site/test_helpers/cms"
	"testing"
	"time"
)

func TestArticlesList(t *testing.T) {
	env := test_helpers.NewTestEnv(nil)
	revert := env.Default()
	defer revert()

	testCases := []struct {
		Name                       string
		Request                    string
		Mock                       cms.MockResponse
		ExpectedResponseJSON       string
		ExpectedResponseStatusCode int
	}{
		{
			Name:    "success",
			Request: "/api/articles?offset=0&limit=20",
			Mock: cms.MockResponse{
				Response: &proto.ArticleListResponse{
					Articles: []*proto.ArticleListItem{
						{
							Title: "Первая",
							Slug:  "pervaya",
							Category: &proto.Category{
								Name: "Быт",
								Slug: "byt",
							},
							Content:     "Short content",
							PublishedAt: timestamppb.New(time.Date(2022, time.May, 24, 19, 54, 20, 0, time.UTC)),
						},
						{
							Title: "Вторая",
							Slug:  "vtoraya",
							Category: &proto.Category{
								Name: "Кухня",
								Slug: "kuhnya",
							},
							Content:     "Short content for second",
							PublishedAt: timestamppb.New(time.Date(2022, time.May, 25, 19, 54, 20, 0, time.UTC)),
						},
					},
					Total: 2,
				},
			},
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON: `{
				"articles": [{
					"title": "Первая",
					"slug": "pervaya",
					"category": {
						"name": "Быт",
						"slug": "byt"
					},
					"content": "Short content",
					"createdAt": "2022-05-24 19:54:20 +0000 UTC"
				}, {
					"title": "Вторая",
					"slug": "vtoraya",
					"category": {
						"name": "Кухня",
						"slug": "kuhnya"
					},
					"content": "Short content for second",
					"createdAt": "2022-05-25 19:54:20 +0000 UTC"
				}],
				"meta": {
					"limit": 20,
					"offset": 0,
					"total": 2
				}
			}`,
		},
	}

	for _, testCase := range testCases {
		var expectedJSON map[string]interface{}
		err := json.Unmarshal([]byte(testCase.ExpectedResponseJSON), &expectedJSON)
		if nil != err {
			t.Fatalf("%s: unmarchal expected json: %s", testCase.Name, err)
		}
		env.CMS.MockForArticleList = testCase.Mock
		env.API.ExecuteTestCases(t, []client.TestCase{
			{
				Message:                    testCase.Name,
				Method:                     http.MethodGet,
				URL:                        testCase.Request,
				ExpectedResponseStatusCode: testCase.ExpectedResponseStatusCode,
				ExpectedResponseJSON:       expectedJSON,
			},
		})
	}
}

func TestArticlesListCategory(t *testing.T) {
	env := test_helpers.NewTestEnv(nil)
	revert := env.Default()
	defer revert()

	testCases := []struct {
		Name                       string
		Request                    string
		Mock                       cms.MockResponse
		ExpectedResponseJSON       string
		ExpectedResponseStatusCode int
	}{
		{
			Name:    "success",
			Request: "/api/articles?offset=0&limit=20&category=test",
			Mock: cms.MockResponse{
				Response: &proto.ArticleListForCategoryResponse{
					Articles: []*proto.ArticleListItem{
						{
							Title: "Первая",
							Slug:  "pervaya",
							Category: &proto.Category{
								Name: "Тест",
								Slug: "test",
							},
							Content:     "Short content",
							PublishedAt: timestamppb.New(time.Date(2022, time.May, 24, 19, 54, 20, 0, time.UTC)),
						},
						{
							Title: "Вторая",
							Slug:  "vtoraya",
							Category: &proto.Category{
								Name: "Тест",
								Slug: "test",
							},
							Content:     "Short content for second",
							PublishedAt: timestamppb.New(time.Date(2022, time.May, 25, 19, 54, 20, 0, time.UTC)),
						},
					},
					Total:        2,
					CategoryName: "Тест",
				},
			},
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON: `{
				"articles": [{
					"title": "Первая",
					"slug": "pervaya",
					"category": {
						"name": "Тест",
						"slug": "test"
					},
					"content": "Short content",
					"createdAt": "2022-05-24 19:54:20 +0000 UTC"
				}, {
					"title": "Вторая",
					"slug": "vtoraya",
					"category": {
						"name": "Тест",
						"slug": "test"
					},
					"content": "Short content for second",
					"createdAt": "2022-05-25 19:54:20 +0000 UTC"
				}],
				"meta": {
					"limit": 20,
					"offset": 0,
					"total": 2,
					"category": "Тест"
				}
			}`,
		},
	}

	for _, testCase := range testCases {
		var expectedJSON map[string]interface{}
		err := json.Unmarshal([]byte(testCase.ExpectedResponseJSON), &expectedJSON)
		if nil != err {
			t.Fatalf("%s: unmarchal expected json: %s", testCase.Name, err)
		}
		env.CMS.MockForArticleListForCategory = testCase.Mock
		env.API.ExecuteTestCases(t, []client.TestCase{
			{
				Message:                    testCase.Name,
				Method:                     http.MethodGet,
				URL:                        testCase.Request,
				ExpectedResponseStatusCode: testCase.ExpectedResponseStatusCode,
				ExpectedResponseJSON:       expectedJSON,
			},
		})
	}
}

func TestArticlesListTag(t *testing.T) {
	env := test_helpers.NewTestEnv(nil)
	revert := env.Default()
	defer revert()

	testCases := []struct {
		Name                       string
		Request                    string
		Mock                       cms.MockResponse
		ExpectedResponseJSON       string
		ExpectedResponseStatusCode int
	}{
		{
			Name:    "success",
			Request: "/api/articles?offset=0&limit=20&tag=test",
			Mock: cms.MockResponse{
				Response: &proto.ArticleListForTagResponse{
					Articles: []*proto.ArticleListItem{
						{
							Title: "Первая",
							Slug:  "pervaya",
							Category: &proto.Category{
								Name: "Кухня",
								Slug: "kuhnya",
							},
							Content:     "Short content",
							PublishedAt: timestamppb.New(time.Date(2022, time.May, 24, 19, 54, 20, 0, time.UTC)),
						},
						{
							Title: "Вторая",
							Slug:  "vtoraya",
							Category: &proto.Category{
								Name: "Быт",
								Slug: "byt",
							},
							Content:     "Short content for second",
							PublishedAt: timestamppb.New(time.Date(2022, time.May, 25, 19, 54, 20, 0, time.UTC)),
						},
					},
					Total:   2,
					TagName: "Тест",
				},
			},
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON: `{
				"articles": [{
					"title": "Первая",
					"slug": "pervaya",
					"category": {
						"name": "Кухня",
						"slug": "kuhnya"
					},
					"content": "Short content",
					"createdAt": "2022-05-24 19:54:20 +0000 UTC"
				}, {
					"title": "Вторая",
					"slug": "vtoraya",
					"category": {
						"name": "Быт",
						"slug": "byt"
					},
					"content": "Short content for second",
					"createdAt": "2022-05-25 19:54:20 +0000 UTC"
				}],
				"meta": {
					"limit": 20,
					"offset": 0,
					"total": 2,
					"tag": "Тест"
				}
			}`,
		},
	}

	for _, testCase := range testCases {
		var expectedJSON map[string]interface{}
		err := json.Unmarshal([]byte(testCase.ExpectedResponseJSON), &expectedJSON)
		if nil != err {
			t.Fatalf("%s: unmarchal expected json: %s", testCase.Name, err)
		}
		env.CMS.MockForArticleListForTag = testCase.Mock
		env.API.ExecuteTestCases(t, []client.TestCase{
			{
				Message:                    testCase.Name,
				Method:                     http.MethodGet,
				URL:                        testCase.Request,
				ExpectedResponseStatusCode: testCase.ExpectedResponseStatusCode,
				ExpectedResponseJSON:       expectedJSON,
			},
		})
	}
}

func TestArticlesListQuery(t *testing.T) {
	env := test_helpers.NewTestEnv(nil)
	revert := env.Default()
	defer revert()

	testCases := []struct {
		Name                       string
		Request                    string
		Mock                       cms.MockResponse
		ExpectedResponseJSON       string
		ExpectedResponseStatusCode int
	}{
		{
			Name:    "success",
			Request: "/api/articles?offset=0&limit=20&q=short",
			Mock: cms.MockResponse{
				Response: &proto.ArticleListResponse{
					Articles: []*proto.ArticleListItem{
						{
							Title: "Первая",
							Slug:  "pervaya",
							Category: &proto.Category{
								Name: "Кухня",
								Slug: "kuhnya",
							},
							Content:     "Short content",
							PublishedAt: timestamppb.New(time.Date(2022, time.May, 24, 19, 54, 20, 0, time.UTC)),
						},
						{
							Title: "Вторая",
							Slug:  "vtoraya",
							Category: &proto.Category{
								Name: "Быт",
								Slug: "byt",
							},
							Content:     "Short content for second",
							PublishedAt: timestamppb.New(time.Date(2022, time.May, 25, 19, 54, 20, 0, time.UTC)),
						},
					},
					Total: 2,
				},
			},
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON: `{
				"articles": [{
					"title": "Первая",
					"slug": "pervaya",
					"category": {
						"name": "Кухня",
						"slug": "kuhnya"
					},
					"content": "Short content",
					"createdAt": "2022-05-24 19:54:20 +0000 UTC"
				}, {
					"title": "Вторая",
					"slug": "vtoraya",
					"category": {
						"name": "Быт",
						"slug": "byt"
					},
					"content": "Short content for second",
					"createdAt": "2022-05-25 19:54:20 +0000 UTC"
				}],
				"meta": {
					"limit": 20,
					"offset": 0,
					"total": 2
				}
			}`,
		},
	}

	for _, testCase := range testCases {
		var expectedJSON map[string]interface{}
		err := json.Unmarshal([]byte(testCase.ExpectedResponseJSON), &expectedJSON)
		if nil != err {
			t.Fatalf("%s: unmarchal expected json: %s", testCase.Name, err)
		}
		env.CMS.MockForArticleListForQuery = testCase.Mock
		env.API.ExecuteTestCases(t, []client.TestCase{
			{
				Message:                    testCase.Name,
				Method:                     http.MethodGet,
				URL:                        testCase.Request,
				ExpectedResponseStatusCode: testCase.ExpectedResponseStatusCode,
				ExpectedResponseJSON:       expectedJSON,
			},
		})
	}
}

func TestArticle(t *testing.T) {
	env := test_helpers.NewTestEnv(nil)
	revert := env.Default()
	defer revert()

	testCases := []struct {
		Name                       string
		Request                    string
		Mock                       cms.MockResponse
		ExpectedResponseJSON       string
		ExpectedResponseStatusCode int
	}{
		{
			Name:    "success",
			Request: "/api/articles/pervaya",
			Mock: cms.MockResponse{
				Response: &proto.ArticleResponse{
					Response: &proto.ArticleResponse_Article{
						Article: &proto.Article{
							Title: "Первая",
							Slug:  "pervaya",
							Category: &proto.Category{
								Name: "Кухня",
								Slug: "kuhnya",
							},
							Content: "Short content",
							Tags: []*proto.Tag{
								{
									Name: "Быт",
									Slug: "byt",
								},
							},
							PublishedAt: timestamppb.New(time.Date(2022, time.May, 24, 19, 54, 20, 0, time.UTC)),
						},
					},
				},
			},
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponseJSON: `{
					"title": "Первая",
					"slug": "pervaya",
					"category": {
						"name": "Кухня",
						"slug": "kuhnya"
					},
					"tags": [ {
						"name": "Быт",
						"slug": "byt"
					} ],
					"content": "Short content",
					"createdAt": "2022-05-24 19:54:20 +0000 UTC"
				}`,
		},
		{
			Name:    "not found",
			Request: "/api/articles/pervaya",
			Mock: cms.MockResponse{
				Response: &proto.ArticleResponse{
					Response: &proto.ArticleResponse_NotFound{NotFound: &proto.NotFound{}}},
			},
			ExpectedResponseStatusCode: http.StatusNotFound,
			ExpectedResponseJSON:       `{"Status":"Not Found"}`,
		},
	}

	for _, testCase := range testCases {
		var expectedJSON map[string]interface{}
		err := json.Unmarshal([]byte(testCase.ExpectedResponseJSON), &expectedJSON)
		if nil != err {
			t.Fatalf("%s: unmarchal expected json: %s", testCase.Name, err)
		}
		env.CMS.MockForArticle = testCase.Mock
		env.API.ExecuteTestCases(t, []client.TestCase{
			{
				Message:                    testCase.Name,
				Method:                     http.MethodGet,
				URL:                        testCase.Request,
				ExpectedResponseStatusCode: testCase.ExpectedResponseStatusCode,
				ExpectedResponseJSON:       expectedJSON,
			},
		})
	}
}
