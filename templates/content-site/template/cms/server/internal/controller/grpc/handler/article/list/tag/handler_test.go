package tag

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
	category := model.Category{
		ID:   "12",
		Slug: "first-category-slug",
		Name: "first-category-name",
	}
	tag := model.Tag{
		ID:   "32",
		Slug: "first-tag-slug",
		Name: "first tag name",
	}
	firstArticle := model.ArticleShortInfoWithCategory{
		ID:          "42",
		Slug:        "first-article-slug",
		Title:       "first article title",
		Category:    category,
		ShortBody:   "first article short body",
		PublishedAt: time.Date(2022, time.January, 1, 23, 12, 4, 0, time.UTC),
	}
	secondArticle := model.ArticleShortInfoWithCategory{
		ID:          "43",
		Slug:        "second-article-slug",
		Title:       "second article title",
		Category:    category,
		ShortBody:   "second article short body",
		PublishedAt: time.Date(2022, time.January, 2, 10, 30, 0, 0, time.UTC),
	}

	testCases := []struct {
		name             string
		tagService       func(t mockConstructorTestingTNewMockTagService) *MockTagService
		articleService   func(t mockConstructorTestingTNewMockArticleService) *MockArticleService
		req              *proto.ArticleListForTagRequest
		expectedResponse *proto.ArticleListForTagResponse
		expectedError    error
	}{
		{
			name: "success",
			tagService: func(t mockConstructorTestingTNewMockTagService) *MockTagService {
				m := NewMockTagService(t)
				m.
					On("GetBySlug", mock.Anything, tag.Slug).
					Return(tag, nil).
					Once()
				return m
			},
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("TotalForTag", mock.Anything, tag.Slug).
					Return(3, nil).
					Once()
				m.
					On(
						"FindShortInfoWithCategoryForTag",
						mock.Anything,
						tag.Slug,
						int64(1),
						int64(2),
					).
					Return([]model.ArticleShortInfoWithCategory{firstArticle, secondArticle}, nil).
					Once()
				return m
			},
			req: &proto.ArticleListForTagRequest{
				TagSlug: string(tag.Slug),
				Offset:  1,
				Limit:   2,
			},
			expectedResponse: &proto.ArticleListForTagResponse{
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
				Total:   3,
				TagName: tag.Name,
			},
		},
		{
			name: "fail get tag",
			tagService: func(t mockConstructorTestingTNewMockTagService) *MockTagService {
				m := NewMockTagService(t)
				m.
					On("GetBySlug", mock.Anything, tag.Slug).
					Return(model.Tag{}, errors.New("test")). //nolint:exhaustruct
					Once()
				return m
			},
			articleService: NewMockArticleService,
			req: &proto.ArticleListForTagRequest{
				TagSlug: string(tag.Slug),
				Offset:  1,
				Limit:   2,
			},
			expectedResponse: nil,
			expectedError:    errors.New("tag get by slug: test"),
		},
		{
			name: "fail get total",
			tagService: func(t mockConstructorTestingTNewMockTagService) *MockTagService {
				m := NewMockTagService(t)
				m.
					On("GetBySlug", mock.Anything, tag.Slug).
					Return(tag, nil).
					Once()
				return m
			},
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("TotalForTag", mock.Anything, tag.Slug).
					Return(0, errors.New("test")).
					Once()
				return m
			},
			req: &proto.ArticleListForTagRequest{
				TagSlug: string(tag.Slug),
				Offset:  1,
				Limit:   2,
			},
			expectedResponse: nil,
			expectedError:    errors.New("total: test"),
		},
		{
			name: "fail get articles",
			tagService: func(t mockConstructorTestingTNewMockTagService) *MockTagService {
				m := NewMockTagService(t)
				m.
					On("GetBySlug", mock.Anything, tag.Slug).
					Return(tag, nil).
					Once()
				return m
			},
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("TotalForTag", mock.Anything, tag.Slug).
					Return(3, nil).
					Once()
				m.
					On(
						"FindShortInfoWithCategoryForTag",
						mock.Anything,
						tag.Slug,
						int64(1),
						int64(2),
					).
					Return([]model.ArticleShortInfoWithCategory{}, errors.New("test")).
					Once()
				return m
			},
			req: &proto.ArticleListForTagRequest{
				TagSlug: string(tag.Slug),
				Offset:  1,
				Limit:   2,
			},
			expectedResponse: nil,
			expectedError:    errors.New("find: test"),
		},
	}

	for _, test := range testCases {
		t.Run(
			test.name, func(t *testing.T) {
				// Arrange
				h := New(test.articleService(t), test.tagService(t))

				// Action
				resp, err := h.ArticleListForTag(context.Background(), test.req)

				// Arrange
				asserts.Equals(t, test.expectedResponse, resp)
				asserts.ErrorsEqual(t, test.expectedError, err)
			},
		)
	}
}
