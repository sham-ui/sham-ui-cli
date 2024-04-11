package detail

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

func TestHandler_Article(t *testing.T) {
	article := model.Article{
		ID:          "42",
		Slug:        "first-article-slug",
		Title:       "first article title",
		CategoryID:  "12",
		ShortBody:   "first article short body",
		Body:        "first article body",
		PublishedAt: time.Date(2022, time.January, 1, 23, 12, 4, 0, time.UTC),
	}
	category := model.Category{
		ID:   "12",
		Slug: "first-category-slug",
		Name: "first category name",
	}
	tags := []model.Tag{
		{
			ID:   "32",
			Slug: "first-tag-slug",
			Name: "first tag name",
		},
		{
			ID:   "33",
			Slug: "second-tag-slug",
			Name: "second tag name",
		},
	}

	testCases := []struct {
		name              string
		req               *proto.ArticleRequest
		articleService    func(t mockConstructorTestingTNewMockArticleService) *MockArticleService
		categoryService   func(t mockConstructorTestingTNewMockCategoryService) *MockCategoryService
		articleTagService func(t mockConstructorTestingTNewMockArticleTagService) *MockArticleTagService
		expectedResponse  *proto.ArticleResponse
		expectedError     error
	}{
		{
			name: "success",
			req: &proto.ArticleRequest{
				Slug: string(article.Slug),
			},
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("FindBySlug", mock.Anything, article.Slug).
					Return(article, nil).
					Once()
				return m
			},
			categoryService: func(t mockConstructorTestingTNewMockCategoryService) *MockCategoryService {
				m := NewMockCategoryService(t)
				m.
					On("GetByID", mock.Anything, article.CategoryID).
					Return(&category, nil).
					Once()
				return m
			},
			articleTagService: func(t mockConstructorTestingTNewMockArticleTagService) *MockArticleTagService {
				m := NewMockArticleTagService(t)
				m.
					On("GetTagForArticle", mock.Anything, article.ID).
					Return(tags, nil).
					Once()
				return m
			},
			expectedResponse: &proto.ArticleResponse{
				Response: &proto.ArticleResponse_Article{
					Article: &proto.Article{
						Slug:         string(article.Slug),
						Title:        article.Title,
						ShortContent: article.ShortBody,
						Content:      article.Body,
						PublishedAt:  timestamppb.New(article.PublishedAt),
						Category: &proto.Category{
							Slug: string(category.Slug),
							Name: category.Name,
						},
						Tags: []*proto.Tag{
							{
								Name: tags[0].Name,
								Slug: string(tags[0].Slug),
							},
							{
								Name: tags[1].Name,
								Slug: string(tags[1].Slug),
							},
						},
					},
				},
			},
		},
		{
			name: "article not found",
			req: &proto.ArticleRequest{
				Slug: string(article.Slug),
			},
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("FindBySlug", mock.Anything, article.Slug).
					Return(model.Article{}, errors.New("article not found")). //nolint:exhaustruct
					Once()
				return m
			},
			categoryService:   NewMockCategoryService,
			articleTagService: NewMockArticleTagService,
			expectedError:     errors.New("find article: article not found"),
		},
		{
			name: "fail get category",
			req: &proto.ArticleRequest{
				Slug: string(article.Slug),
			},
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("FindBySlug", mock.Anything, article.Slug).
					Return(article, nil).
					Once()
				return m
			},
			categoryService: func(t mockConstructorTestingTNewMockCategoryService) *MockCategoryService {
				m := NewMockCategoryService(t)
				m.
					On("GetByID", mock.Anything, article.CategoryID).
					Return(nil, errors.New("test")).
					Once()
				return m
			},
			articleTagService: NewMockArticleTagService,
			expectedError:     errors.New("category get by id: test"),
		},
		{
			name: "fail get tag",
			req: &proto.ArticleRequest{
				Slug: string(article.Slug),
			},
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("FindBySlug", mock.Anything, article.Slug).
					Return(article, nil).
					Once()
				return m
			},
			categoryService: func(t mockConstructorTestingTNewMockCategoryService) *MockCategoryService {
				m := NewMockCategoryService(t)
				m.
					On("GetByID", mock.Anything, article.CategoryID).
					Return(&category, nil).
					Once()
				return m
			},
			articleTagService: func(t mockConstructorTestingTNewMockArticleTagService) *MockArticleTagService {
				m := NewMockArticleTagService(t)
				m.
					On("GetTagForArticle", mock.Anything, article.ID).
					Return([]model.Tag{}, errors.New("test")).
					Once()
				return m
			},
			expectedError: errors.New("get tags for article: test"),
		},
	}

	for _, test := range testCases {
		t.Run(
			test.name, func(t *testing.T) {
				// Arrange
				articleService := test.articleService(t)
				categoryService := test.categoryService(t)
				tagService := test.articleTagService(t)
				h := NewHandler(articleService, categoryService, tagService)

				// Action
				resp, err := h.Article(context.Background(), test.req)

				// Asserts
				asserts.Equals(t, test.expectedResponse, resp)
				asserts.ErrorsEqual(t, test.expectedError, err)
			},
		)
	}
}
