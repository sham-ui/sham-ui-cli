package detault

import (
	"cms/internal/controller/grpc/proto"
	"cms/internal/model"
	"cms/pkg/asserts"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestHandler_ArticleList(t *testing.T) {
	firstArticle := model.ArticleShortInfoWithCategory{
		ID:    "42",
		Slug:  "first-article-slug",
		Title: "first article title",
		Category: model.Category{
			ID:   "12",
			Slug: "first-category-slug",
			Name: "first-category-name",
		},
		ShortBody:   "first article short body",
		PublishedAt: time.Date(2022, time.January, 1, 23, 12, 4, 0, time.UTC),
	}
	secondArticle := model.ArticleShortInfoWithCategory{
		ID:    "43",
		Slug:  "second-article-slug",
		Title: "second article title",
		Category: model.Category{
			ID:   "13",
			Slug: "second-category-slug",
			Name: "second-category-name",
		},
		ShortBody:   "second article short body",
		PublishedAt: time.Date(2022, time.January, 2, 10, 30, 0, 0, time.UTC),
	}

	testCases := []struct {
		name             string
		articleService   func(t mockConstructorTestingTNewMockArticleService) *MockArticleService
		req              *proto.ArticleListRequest
		expectedResponse *proto.ArticleListResponse
		expectedError    error
	}{
		{
			name: "success",
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("Total", mock.Anything).
					Return(3, nil).
					Once()
				m.
					On("FindShortInfoWithCategory", mock.Anything, int64(1), int64(2)).
					Return([]model.ArticleShortInfoWithCategory{firstArticle, secondArticle}, nil).
					Once()
				return m
			},
			req: &proto.ArticleListRequest{Offset: 1, Limit: 2},
			expectedResponse: &proto.ArticleListResponse{
				Articles: []*proto.ArticleListItem{
					{
						Title: firstArticle.Title,
						Slug:  string(firstArticle.Slug),
						Category: &proto.Category{
							Name: firstArticle.Category.Name,
							Slug: string(firstArticle.Category.Slug),
						},
						Content:     firstArticle.ShortBody,
						PublishedAt: timestamppb.New(firstArticle.PublishedAt),
					},
					{
						Title: secondArticle.Title,
						Slug:  string(secondArticle.Slug),
						Category: &proto.Category{
							Name: secondArticle.Category.Name,
							Slug: string(secondArticle.Category.Slug),
						},
						Content:     secondArticle.ShortBody,
						PublishedAt: timestamppb.New(secondArticle.PublishedAt),
					},
				},
				Total: 3,
			},
		},
		{
			name: "fail get total",
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("Total", mock.Anything).
					Return(0, errors.New("test")).
					Once()
				return m
			},
			req:              &proto.ArticleListRequest{Offset: 1, Limit: 2},
			expectedResponse: nil,
			expectedError:    errors.New("total: test"),
		},
		{
			name: "fail get articles",
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("Total", mock.Anything).
					Return(3, nil).
					Once()
				m.
					On("FindShortInfoWithCategory", mock.Anything, int64(1), int64(2)).
					Return([]model.ArticleShortInfoWithCategory{}, errors.New("test")).
					Once()
				return m
			},
			req:              &proto.ArticleListRequest{Offset: 1, Limit: 2},
			expectedResponse: nil,
			expectedError:    errors.New("find: test"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			h := New(test.articleService(t))

			// Action
			resp, err := h.ArticleList(context.Background(), test.req)

			// Arrange
			asserts.Equals(t, test.expectedResponse, resp)
			asserts.ErrorsEqual(t, test.expectedError, err)
		})
	}
}
