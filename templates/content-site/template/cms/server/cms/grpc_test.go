package main

import (
	"cms/api"
	"cms/proto"
	"cms/test_helpers"
	"cms/test_helpers/asserts"
	"cms/test_helpers/database"
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestGRPC(t *testing.T) {
	publishedAt := time.Date(2022, time.March, 16, 18, 45, 0, 0, time.UTC)

	clearArticles := func(execSQL database.SQLQueryExecutor) {
		execSQL("DELETE FROM article_tag")
		execSQL("DELETE FROM article")
		execSQL("DELETE FROM category")
		execSQL("DELETE FROM tag")
	}

	createCategory := func(execSQL database.SQLQueryExecutor) {
		execSQL("INSERT into category(id, name, slug) VALUES(1, 'Хобби', 'hobby')")
	}

	createTags := func(execSQL database.SQLQueryExecutor) {
		execSQL("INSERT into tag(id, name, slug) VALUES(1, 'Первый', 'first')")
		execSQL("INSERT into tag(id, name, slug) VALUES(2, 'Второй', 'second')")
	}

	createArticle := func(execSQL database.SQLQueryExecutor) {
		createCategory(execSQL)
		createTags(execSQL)
		execSQL(
			"INSERT INTO "+
				"article(id, title, slug, category_id, short_body, body, published_at) "+
				"VALUES($1, $2, $3, $4, $5, $6, $7)",
			1,
			"test slug",
			"test-slug",
			1,
			"short body",
			"body text",
			publishedAt,
		)
	}
	createArticleWithTag := func(execSQL database.SQLQueryExecutor) {
		createArticle(execSQL)
		execSQL("INSERT INTO article_tag(id, article_id, tag_id) VALUES(1, 1, 1)")
		execSQL("INSERT INTO article_tag(id, article_id, tag_id) VALUES(2, 1, 2)")
	}

	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()

	srv := api.NewAPI(env.DB.DB)

	t.Run("ArticleList", func(t *testing.T) {
		testCases := []struct {
			Name             string
			Prepare          func(execSQL database.SQLQueryExecutor)
			Request          *proto.ArticleListRequest
			ExpectedResponse *proto.ArticleListResponse
			ExpectedError    error
		}{
			{
				Name: "Empty",
				Request: &proto.ArticleListRequest{
					Limit: 10,
				},
				ExpectedResponse: &proto.ArticleListResponse{},
			},
			{
				Name:    "One item",
				Prepare: createArticle,
				Request: &proto.ArticleListRequest{
					Limit: 10,
				},
				ExpectedResponse: &proto.ArticleListResponse{
					Articles: []*proto.ArticleListItem{
						{
							Title: "test slug",
							Slug:  "test-slug",
							Category: &proto.Category{
								Name: "Хобби",
								Slug: "hobby",
							},
							Content:     "short body",
							PublishedAt: timestamppb.New(publishedAt),
						},
					},
					Total: 1,
				},
			},
			{
				Name:    "limit == 0",
				Prepare: createArticle,
				Request: &proto.ArticleListRequest{},
				ExpectedResponse: &proto.ArticleListResponse{
					Total: 1,
				},
			},
			{
				Name:    "offset == 1",
				Prepare: createArticle,
				Request: &proto.ArticleListRequest{
					Offset: 1,
					Limit:  10,
				},
				ExpectedResponse: &proto.ArticleListResponse{
					Total: 1,
				},
			},
		}
		for _, testCase := range testCases {
			execSQL := env.DB.ExecForCase(t, testCase.Name)
			clearArticles(execSQL)
			if nil != testCase.Prepare {
				testCase.Prepare(execSQL)
			}
			resp, err := srv.ArticleList(context.Background(), testCase.Request)
			asserts.Equals(t, testCase.ExpectedError, err, fmt.Sprintf("%s: check error", testCase.Name))
			asserts.Equals(t, testCase.ExpectedResponse, resp, fmt.Sprintf("%s: check response", testCase.Name))
		}
	})

	t.Run("ArticleListForCategory", func(t *testing.T) {
		testCases := []struct {
			Name             string
			Prepare          func(execSQL database.SQLQueryExecutor)
			Request          *proto.ArticleListForCategoryRequest
			ExpectedResponse *proto.ArticleListForCategoryResponse
			ExpectedError    error
		}{
			{
				Name:    "Empty",
				Prepare: createCategory,
				Request: &proto.ArticleListForCategoryRequest{
					CategorySlug: "hobby",
					Limit:        10,
				},
				ExpectedResponse: &proto.ArticleListForCategoryResponse{
					CategoryName: "Хобби",
				},
			},
			{
				Name:    "One item",
				Prepare: createArticle,
				Request: &proto.ArticleListForCategoryRequest{
					CategorySlug: "hobby",
					Limit:        10,
				},
				ExpectedResponse: &proto.ArticleListForCategoryResponse{
					Articles: []*proto.ArticleListItem{
						{
							Title: "test slug",
							Slug:  "test-slug",
							Category: &proto.Category{
								Name: "Хобби",
								Slug: "hobby",
							},
							Content:     "short body",
							PublishedAt: timestamppb.New(publishedAt),
						},
					},
					Total:        1,
					CategoryName: "Хобби",
				},
			},
			{
				Name:    "limit == 0",
				Prepare: createArticle,
				Request: &proto.ArticleListForCategoryRequest{
					CategorySlug: "hobby",
				},
				ExpectedResponse: &proto.ArticleListForCategoryResponse{
					Total:        1,
					CategoryName: "Хобби",
				},
			},
			{
				Name:    "offset == 1",
				Prepare: createArticle,
				Request: &proto.ArticleListForCategoryRequest{
					CategorySlug: "hobby",
					Offset:       1,
					Limit:        10,
				},
				ExpectedResponse: &proto.ArticleListForCategoryResponse{
					Total:        1,
					CategoryName: "Хобби",
				},
			},
			{
				Name:    "not found",
				Prepare: createArticle,
				Request: &proto.ArticleListForCategoryRequest{
					CategorySlug: "hobby-not-found",
					Limit:        10,
				},
				ExpectedResponse: nil,
				ExpectedError:    errors.New("category: select string: sql: no rows in result set"),
			},
		}
		for _, testCase := range testCases {
			execSQL := env.DB.ExecForCase(t, testCase.Name)
			clearArticles(execSQL)
			if nil != testCase.Prepare {
				testCase.Prepare(execSQL)
			}
			resp, err := srv.ArticleListForCategory(context.Background(), testCase.Request)
			asserts.Equals(t, testCase.ExpectedError, err, fmt.Sprintf("%s: check error", testCase.Name))
			asserts.Equals(t, testCase.ExpectedResponse, resp, fmt.Sprintf("%s: check response", testCase.Name))
		}
	})

	t.Run("ArticleListForTag", func(t *testing.T) {
		testCases := []struct {
			Name             string
			Prepare          func(execSQL database.SQLQueryExecutor)
			Request          *proto.ArticleListForTagRequest
			ExpectedResponse *proto.ArticleListForTagResponse
			ExpectedError    error
		}{
			{
				Name:    "Empty",
				Prepare: createTags,
				Request: &proto.ArticleListForTagRequest{
					TagSlug: "first",
					Limit:   10,
				},
				ExpectedResponse: &proto.ArticleListForTagResponse{
					TagName: "Первый",
				},
			},
			{
				Name:    "One item",
				Prepare: createArticleWithTag,
				Request: &proto.ArticleListForTagRequest{
					TagSlug: "first",
					Limit:   10,
				},
				ExpectedResponse: &proto.ArticleListForTagResponse{
					Articles: []*proto.ArticleListItem{
						{
							Title: "test slug",
							Slug:  "test-slug",
							Category: &proto.Category{
								Name: "Хобби",
								Slug: "hobby",
							},
							Content:     "short body",
							PublishedAt: timestamppb.New(publishedAt),
						},
					},
					Total:   1,
					TagName: "Первый",
				},
			},
			{
				Name:    "limit == 0",
				Prepare: createArticleWithTag,
				Request: &proto.ArticleListForTagRequest{
					TagSlug: "first",
				},
				ExpectedResponse: &proto.ArticleListForTagResponse{
					Total:   1,
					TagName: "Первый",
				},
			},
			{
				Name:    "offset == 1",
				Prepare: createArticleWithTag,
				Request: &proto.ArticleListForTagRequest{
					TagSlug: "first",
					Offset:  1,
					Limit:   10,
				},
				ExpectedResponse: &proto.ArticleListForTagResponse{
					Total:   1,
					TagName: "Первый",
				},
			},
			{
				Name:    "not found",
				Prepare: createArticleWithTag,
				Request: &proto.ArticleListForTagRequest{
					TagSlug: "third-not-found",
					Limit:   10,
				},
				ExpectedResponse: nil,
				ExpectedError:    errors.New("tag: select string: sql: no rows in result set"),
			},
		}
		for _, testCase := range testCases {
			execSQL := env.DB.ExecForCase(t, testCase.Name)
			clearArticles(execSQL)
			if nil != testCase.Prepare {
				testCase.Prepare(execSQL)
			}
			resp, err := srv.ArticleListForTag(context.Background(), testCase.Request)
			asserts.Equals(t, testCase.ExpectedError, err, fmt.Sprintf("%s: check error", testCase.Name))
			asserts.Equals(t, testCase.ExpectedResponse, resp, fmt.Sprintf("%s: check response", testCase.Name))
		}
	})

	t.Run("ArticleListForQuery", func(t *testing.T) {
		testCases := []struct {
			Name             string
			Prepare          func(execSQL database.SQLQueryExecutor)
			Request          *proto.ArticleListForQueryRequest
			ExpectedResponse *proto.ArticleListResponse
			ExpectedError    error
		}{
			{
				Name: "Empty",
				Request: &proto.ArticleListForQueryRequest{
					Query: "first",
					Limit: 10,
				},
				ExpectedResponse: &proto.ArticleListResponse{},
			},
			{
				Name:    "One item (title)",
				Prepare: createArticle,
				Request: &proto.ArticleListForQueryRequest{
					Query: "test",
					Limit: 10,
				},
				ExpectedResponse: &proto.ArticleListResponse{
					Articles: []*proto.ArticleListItem{
						{
							Title: "test slug",
							Slug:  "test-slug",
							Category: &proto.Category{
								Name: "Хобби",
								Slug: "hobby",
							},
							Content:     "short body",
							PublishedAt: timestamppb.New(publishedAt),
						},
					},
					Total: 1,
				},
			},
			{
				Name:    "One item (short_body)",
				Prepare: createArticle,
				Request: &proto.ArticleListForQueryRequest{
					Query: "short",
					Limit: 10,
				},
				ExpectedResponse: &proto.ArticleListResponse{
					Articles: []*proto.ArticleListItem{
						{
							Title: "test slug",
							Slug:  "test-slug",
							Category: &proto.Category{
								Name: "Хобби",
								Slug: "hobby",
							},
							Content:     "short body",
							PublishedAt: timestamppb.New(publishedAt),
						},
					},
					Total: 1,
				},
			},
			{
				Name:    "One item (body)",
				Prepare: createArticle,
				Request: &proto.ArticleListForQueryRequest{
					Query: "body",
					Limit: 10,
				},
				ExpectedResponse: &proto.ArticleListResponse{
					Articles: []*proto.ArticleListItem{
						{
							Title: "test slug",
							Slug:  "test-slug",
							Category: &proto.Category{
								Name: "Хобби",
								Slug: "hobby",
							},
							Content:     "short body",
							PublishedAt: timestamppb.New(publishedAt),
						},
					},
					Total: 1,
				},
			},
			{
				Name:    "limit == 0",
				Prepare: createArticle,
				Request: &proto.ArticleListForQueryRequest{
					Query: "test",
				},
				ExpectedResponse: &proto.ArticleListResponse{
					Total: 1,
				},
			},
			{
				Name:    "offset == 1",
				Prepare: createArticle,
				Request: &proto.ArticleListForQueryRequest{
					Query:  "test",
					Offset: 1,
					Limit:  10,
				},
				ExpectedResponse: &proto.ArticleListResponse{
					Total: 1,
				},
			},
			{
				Name:    "not found",
				Prepare: createArticle,
				Request: &proto.ArticleListForQueryRequest{
					Query: "third-not-found",
					Limit: 10,
				},
				ExpectedResponse: &proto.ArticleListResponse{
					Total: 0,
				},
			},
		}
		for _, testCase := range testCases {
			execSQL := env.DB.ExecForCase(t, testCase.Name)
			clearArticles(execSQL)
			if nil != testCase.Prepare {
				testCase.Prepare(execSQL)
			}
			resp, err := srv.ArticleListForQuery(context.Background(), testCase.Request)
			asserts.Equals(t, testCase.ExpectedError, err, fmt.Sprintf("%s: check error", testCase.Name))
			asserts.Equals(t, testCase.ExpectedResponse, resp, fmt.Sprintf("%s: check response", testCase.Name))
		}
	})

	t.Run("Article", func(t *testing.T) {
		testCases := []struct {
			Name             string
			Prepare          func(execSQL database.SQLQueryExecutor)
			Request          *proto.ArticleRequest
			ExpectedResponse *proto.ArticleResponse
			ExpectedError    error
		}{
			{
				Name:    "Empty",
				Request: &proto.ArticleRequest{},
				ExpectedResponse: &proto.ArticleResponse{
					Response: &proto.ArticleResponse_NotFound{
						NotFound: &proto.NotFound{},
					},
				},
			},
			{
				Name:    "Empty slug",
				Prepare: createArticleWithTag,
				Request: &proto.ArticleRequest{},
				ExpectedResponse: &proto.ArticleResponse{
					Response: &proto.ArticleResponse_NotFound{
						NotFound: &proto.NotFound{},
					},
				},
			},
			{
				Name:    "Not found",
				Prepare: createArticleWithTag,
				Request: &proto.ArticleRequest{
					Slug: "not-found",
				},
				ExpectedResponse: &proto.ArticleResponse{
					Response: &proto.ArticleResponse_NotFound{
						NotFound: &proto.NotFound{},
					},
				},
			},
			{
				Name:    "test-slug",
				Prepare: createArticleWithTag,
				Request: &proto.ArticleRequest{
					Slug: "test-slug",
				},
				ExpectedResponse: &proto.ArticleResponse{
					Response: &proto.ArticleResponse_Article{
						Article: &proto.Article{
							Title: "test slug",
							Slug:  "test-slug",
							Category: &proto.Category{
								Name: "Хобби",
								Slug: "hobby",
							},
							Content: "body text",
							Tags: []*proto.Tag{
								{
									Name: "Первый",
									Slug: "first",
								},
								{
									Name: "Второй",
									Slug: "second",
								},
							},
							PublishedAt: timestamppb.New(publishedAt),
						},
					},
				},
			},
		}
		for _, testCase := range testCases {
			execSQL := env.DB.ExecForCase(t, testCase.Name)
			clearArticles(execSQL)
			if nil != testCase.Prepare {
				testCase.Prepare(execSQL)
			}
			resp, err := srv.Article(context.Background(), testCase.Request)
			asserts.Equals(t, testCase.ExpectedError, err, fmt.Sprintf("%s: check error", testCase.Name))
			asserts.Equals(t, testCase.ExpectedResponse, resp, fmt.Sprintf("%s: check response", testCase.Name))
		}
	})
}
