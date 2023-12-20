package cms

import (
	"context"
	"errors"
	"testing"
	"time"

	"site/internal/external_api/cms/proto"
	"site/internal/model"
	"site/pkg/asserts"

	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var protoArticles = []*proto.ArticleListItem{
	{
		Title: "first title",
		Slug:  "first-slug",
		Category: &proto.Category{
			Name: "first-category-name",
			Slug: "first-category-slug",
		},
		Content:     "first-content",
		PublishedAt: timestamppb.New(time.Date(2023, time.November, 16, 20, 15, 17, 123, time.UTC)),
	},
	{
		Title: "second title",
		Slug:  "second-slug",
		Category: &proto.Category{
			Name: "second-category-name",
			Slug: "second-category-slug",
		},
		Content:     "second-content",
		PublishedAt: timestamppb.New(time.Date(2022, time.June, 10, 13, 42, 34, 546, time.UTC)),
	},
}

var articles = []model.ShortArticle{
	{
		Title: "first title",
		Slug:  "first-slug",
		Category: model.Category{
			Name: "first-category-name",
			Slug: "first-category-slug",
		},
		ShortContent: "first-content",
		PublishedAt:  time.Date(2023, time.November, 16, 20, 15, 17, 123, time.UTC),
	},
	{
		Title: "second title",
		Slug:  "second-slug",
		Category: model.Category{
			Name: "second-category-name",
			Slug: "second-category-slug",
		},
		ShortContent: "second-content",
		PublishedAt:  time.Date(2022, time.June, 10, 13, 42, 34, 546, time.UTC),
	},
}

func TestClient_Asset(t *testing.T) {
	testCases := []struct {
		name             string
		client           func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient
		expectedResponse *model.Asset
		expectedError    error
	}{
		{
			name: "success",
			client: func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient {
				m := newMockCmsClient(t)
				m.
					On("Asset", mock.Anything, &proto.AssetRequest{Path: "/test.txt"}).
					Return(&proto.AssetResponse{
						Response: &proto.AssetResponse_File{
							File: []byte("content"),
						},
					}, nil).
					Once()
				return m
			},
			expectedResponse: &model.Asset{
				Path:    "/test.txt",
				Content: []byte("content"),
			},
		},
		{
			name: "not found",
			client: func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient {
				m := newMockCmsClient(t)
				m.
					On("Asset", mock.Anything, &proto.AssetRequest{Path: "/test.txt"}).
					Return(&proto.AssetResponse{
						Response: &proto.AssetResponse_NotFound{NotFound: &proto.NotFound{}},
					}, nil).
					Once()
				return m
			},
			expectedError: model.NewAssetNotFoundError("/test.txt"),
		},
		{
			name: "error",
			client: func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient {
				m := newMockCmsClient(t)
				m.
					On("Asset", mock.Anything, &proto.AssetRequest{Path: "/test.txt"}).
					Return(nil, errors.New("test error")).
					Once()
				return m
			},
			expectedError: errors.New("grpc error: Asset: test error"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			resp, err := (&client{pb: test.client(t)}).Asset(context.Background(), "/test.txt")
			asserts.ErrorsEqual(t, test.expectedError, err)
			asserts.Equals(t, test.expectedResponse, resp)
		})
	}
}

func TestClient_Article(t *testing.T) {
	testCases := []struct {
		name             string
		client           func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient
		expectedResponse *model.Article
		expectedError    error
	}{
		{
			name: "success",
			client: func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient {
				m := newMockCmsClient(t)
				m.
					On("Article", mock.Anything, &proto.ArticleRequest{Slug: "slug"}).
					Return(&proto.ArticleResponse{
						Response: &proto.ArticleResponse_Article{
							Article: &proto.Article{
								Title: "title",
								Slug:  "slug",
								Category: &proto.Category{
									Name: "category-name",
									Slug: "category-slug",
								},
								ShortContent: "short content",
								Content:      "content",
								Tags: []*proto.Tag\{{
									Name: "tag-name",
									Slug: "tag-slug",
								}},
								PublishedAt: timestamppb.New(time.Date(2023, time.November, 16, 20, 15, 17, 123, time.UTC)),
							},
						},
					}, nil).
					Once()
				return m
			},
			expectedResponse: &model.Article{
				ShortArticle: model.ShortArticle{
					Title: "title",
					Slug:  "slug",
					Category: model.Category{
						Name: "category-name",
						Slug: "category-slug",
					},
					ShortContent: "short content",
					PublishedAt:  time.Date(2023, time.November, 16, 20, 15, 17, 123, time.UTC),
				},
				Tags: []model.Tag\{{
					Name: "tag-name",
					Slug: "tag-slug",
				}},
				Content: "content",
			},
		},
		{
			name: "not found",
			client: func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient {
				m := newMockCmsClient(t)
				m.
					On("Article", mock.Anything, &proto.ArticleRequest{Slug: "slug"}).
					Return(&proto.ArticleResponse{
						Response: &proto.ArticleResponse_NotFound{NotFound: &proto.NotFound{}},
					}, nil).
					Once()
				return m
			},
			expectedError: model.NewArticleNotFoundError("slug"),
		},
		{
			name: "error",
			client: func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient {
				m := newMockCmsClient(t)
				m.
					On("Article", mock.Anything, &proto.ArticleRequest{Slug: "slug"}).
					Return(nil, errors.New("test error")).
					Once()
				return m
			},
			expectedError: errors.New("grpc error: Article: test error"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			resp, err := (&client{pb: test.client(t)}).Article(context.Background(), "slug")
			asserts.ErrorsEqual(t, test.expectedError, err)
			asserts.Equals(t, test.expectedResponse, resp)
		})
	}
}

func TestClient_Articles(t *testing.T) {
	testCases := []struct {
		name             string
		client           func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient
		expectedResponse *model.PaginatedArticles
		expectedError    error
	}{
		{
			name: "success",
			client: func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient {
				m := newMockCmsClient(t)
				m.
					On("ArticleList", mock.Anything, &proto.ArticleListRequest{Limit: 30, Offset: 10}).
					Return(&proto.ArticleListResponse{
						Articles: protoArticles,
						Total:    35,
					}, nil).
					Once()
				return m
			},
			expectedResponse: &model.PaginatedArticles{
				Articles: articles,
				Total:    35,
			},
		},
		{
			name: "empty",
			client: func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient {
				m := newMockCmsClient(t)
				m.
					On("ArticleList", mock.Anything, &proto.ArticleListRequest{Limit: 30, Offset: 10}).
					Return(&proto.ArticleListResponse{
						Articles: []*proto.ArticleListItem{},
						Total:    0,
					}, nil).
					Once()
				return m
			},
			expectedResponse: &model.PaginatedArticles{
				Articles: []model.ShortArticle{},
				Total:    0,
			},
		},
		{
			name: "error",
			client: func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient {
				m := newMockCmsClient(t)
				m.
					On("ArticleList", mock.Anything, &proto.ArticleListRequest{Limit: 30, Offset: 10}).
					Return(nil, errors.New("test error")).
					Once()
				return m
			},
			expectedError: errors.New("grpc error: ArticleList: test error"),
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			resp, err := (&client{pb: test.client(t)}).Articles(context.Background(), 10, 30)
			asserts.ErrorsEqual(t, test.expectedError, err)
			asserts.Equals(t, test.expectedResponse, resp)
		})
	}
}

func TestClient_ArticleListForQuery(t *testing.T) {
	testCases := []struct {
		name             string
		client           func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient
		expectedResponse *model.PaginatedArticles
		expectedError    error
	}{
		{
			name: "success",
			client: func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient {
				m := newMockCmsClient(t)
				m.
					On("ArticleListForQuery", mock.Anything, &proto.ArticleListForQueryRequest{
						Offset: 10,
						Limit:  30,
						Query:  "query",
					}).
					Return(&proto.ArticleListResponse{
						Articles: protoArticles,
						Total:    35,
					}, nil).
					Once()
				return m
			},
			expectedResponse: &model.PaginatedArticles{
				Articles: articles,
				Total:    35,
			},
		},
		{
			name: "error",
			client: func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient {
				m := newMockCmsClient(t)
				m.
					On("ArticleListForQuery", mock.Anything, &proto.ArticleListForQueryRequest{
						Offset: 10,
						Limit:  30,
						Query:  "query",
					}).
					Return(nil, errors.New("test error")).
					Once()
				return m
			},
			expectedError: errors.New("grpc error: ArticleListForQuery: test error"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			resp, err := (&client{pb: test.client(t)}).ArticleListForQuery(context.Background(), "query", 10, 30)
			asserts.ErrorsEqual(t, test.expectedError, err)
			asserts.Equals(t, test.expectedResponse, resp)
		})
	}
}

//nolint:dupl
func TestClient_ArticleListForTag(t *testing.T) {
	testCases := []struct {
		name             string
		client           func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient
		expectedResponse *model.PaginatedArticleForTag
		expectedError    error
	}{
		{
			name: "success",
			client: func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient {
				m := newMockCmsClient(t)
				m.
					On("ArticleListForTag", mock.Anything, &proto.ArticleListForTagRequest{
						TagSlug: "tag-slug",
						Offset:  10,
						Limit:   30,
					}).
					Return(&proto.ArticleListForTagResponse{
						Articles: protoArticles,
						TagName:  "tag-name",
						Total:    35,
					}, nil).
					Once()
				return m
			},
			expectedResponse: &model.PaginatedArticleForTag{
				Articles: articles,
				TagName:  "tag-name",
				Total:    35,
			},
		},
		{
			name: "error",
			client: func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient {
				m := newMockCmsClient(t)
				m.
					On("ArticleListForTag", mock.Anything, &proto.ArticleListForTagRequest{
						TagSlug: "tag-slug",
						Offset:  10,
						Limit:   30,
					}).
					Return(nil, errors.New("test error")).
					Once()
				return m
			},
			expectedError: errors.New("grpc error: ArticleListForTag: test error"),
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			resp, err := (&client{pb: test.client(t)}).ArticleListForTag(context.Background(), "tag-slug", 10, 30)
			asserts.ErrorsEqual(t, test.expectedError, err)
			asserts.Equals(t, test.expectedResponse, resp)
		})
	}
}

//nolint:dupl
func TestClient_ArticleListForCategory(t *testing.T) {
	testCases := []struct {
		name             string
		client           func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient
		expectedResponse *model.PaginatedArticleForCategory
		expectedError    error
	}{
		{
			name: "success",
			client: func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient {
				m := newMockCmsClient(t)
				m.
					On("ArticleListForCategory", mock.Anything, &proto.ArticleListForCategoryRequest{
						CategorySlug: "category-slug",
						Offset:       10,
						Limit:        30,
					}).
					Return(&proto.ArticleListForCategoryResponse{
						Articles:     protoArticles,
						Total:        35,
						CategoryName: "category-name",
					}, nil).
					Once()
				return m
			},
			expectedResponse: &model.PaginatedArticleForCategory{
				Articles:     articles,
				Total:        35,
				CategoryName: "category-name",
			},
		},
		{
			name: "error",
			client: func(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient {
				m := newMockCmsClient(t)
				m.
					On("ArticleListForCategory", mock.Anything, &proto.ArticleListForCategoryRequest{
						CategorySlug: "category-slug",
						Offset:       10,
						Limit:        30,
					}).
					Return(nil, errors.New("test error")).
					Once()
				return m
			},
			expectedError: errors.New("grpc error: ArticleListForCategory: test error"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			resp, err := (&client{pb: test.client(t)}).ArticleListForCategory(context.Background(), "category-slug", 10, 30)
			asserts.ErrorsEqual(t, test.expectedError, err)
			asserts.Equals(t, test.expectedResponse, resp)
		})
	}
}
