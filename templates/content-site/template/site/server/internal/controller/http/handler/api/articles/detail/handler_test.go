package detail

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"site/internal/model"
	"site/pkg/asserts"
	"site/pkg/logger"

	"github.com/go-logr/logr/testr"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ServeHTTP(t *testing.T) {
	testCases := []struct {
		name         string
		service      func(t mockConstructorTestingTNewMockArticlesService) *MockArticlesService
		request      *http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name: "success",
			service: func(t mockConstructorTestingTNewMockArticlesService) *MockArticlesService {
				m := NewMockArticlesService(t)
				m.On("Article", mock.Anything, "test-article").Return(&model.Article{
					ShortArticle: model.ShortArticle{
						Title: "title value",
						Slug:  "test-article",
						Category: model.Category{
							Name: "category-value",
							Slug: "category-slug",
						},
						ShortContent: "short content",
						PublishedAt:  time.Date(2022, time.April, 10, 13, 32, 10, 11, time.UTC),
					},
					Tags: []model.Tag{
						{
							Name: "foo",
							Slug: "foo",
						},
						{
							Name: "bar",
							Slug: "bar",
						},
					},
					Content: "content",
				}, nil)
				return m
			},
			request:      httptest.NewRequest(http.MethodGet, "/articles/test-article", nil),
			expectedCode: http.StatusOK,
			expectedBody: `{
			  "title": "title value",
			  "slug": "test-article",
			  "category": {
				"name": "category-value",
				"slug": "category-slug"
			  },
			  "tags": [
				{
				  "name": "foo",
				  "slug": "foo"
				},
				{
				  "name": "bar",
				  "slug": "bar"
				}
			  ],
			  "shortContent": "short content",
			  "content": "content",
			  "createdAt": "2022-04-10 13:32:10.000000011 +0000 UTC"
			}`,
		},
		{
			name: "not found",
			service: func(t mockConstructorTestingTNewMockArticlesService) *MockArticlesService {
				m := NewMockArticlesService(t)
				m.On("Article", mock.Anything, "test-article").Return(nil, model.ArticleNotFoundError{})
				return m
			},
			request:      httptest.NewRequest(http.MethodGet, "/articles/test-article", nil),
			expectedCode: http.StatusNotFound,
			expectedBody: `{"Status":"Not Found", "Messages": ["article not found"]}`,
		},
		{
			name: "internal error",
			service: func(t mockConstructorTestingTNewMockArticlesService) *MockArticlesService {
				m := NewMockArticlesService(t)
				m.On("Article", mock.Anything, "test-article").Return(nil, errors.New("internal error"))
				return m
			},
			request:      httptest.NewRequest(http.MethodGet, "/articles/test-article", nil),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"Status":"Internal Server Error", "Messages": ["internal server error"]}`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := mux.NewRouter()
			Setup(router, test.service(t))
			ctx := logger.Save(test.request.Context(), testr.New(t))
			request := test.request.WithContext(ctx)
			resp := httptest.NewRecorder()

			// Act
			router.ServeHTTP(resp, request)

			// Assert
			asserts.JSONEquals(t, test.expectedBody, resp.Body.String(), "body")
			asserts.Equals(t, test.expectedCode, resp.Code, "code")
		})
	}
}
