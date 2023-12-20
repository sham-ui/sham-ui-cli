package query

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
				m.On("ArticleListForQuery", mock.Anything, "test-query", int64(10), int64(30)).Return(&model.PaginatedArticles{
					Total: 1,
					Articles: []model.ShortArticle{
						{
							Title: "first title",
							Slug:  "first-title",
							Category: model.Category{
								Name: "first-category",
								Slug: "first-category-slug",
							},
							ShortContent: "first short content",
							PublishedAt:  time.Date(2022, time.April, 10, 13, 32, 10, 11, time.UTC),
						},
						{
							Title: "second title",
							Slug:  "second-title",
							Category: model.Category{
								Name: "second-category",
								Slug: "second-category-slug",
							},
							ShortContent: "second short content",
							PublishedAt:  time.Date(2023, time.March, 5, 20, 16, 30, 20, time.UTC),
						},
					},
				}, nil)
				return m
			},
			request:      httptest.NewRequest(http.MethodGet, "/articles?offset=10&limit=30&q=test-query", nil),
			expectedCode: http.StatusOK,
			expectedBody: `{
			  "articles": [
				{
				  "title": "first title",
				  "slug": "first-title",
				  "category": {
					"name": "first-category",
					"slug": "first-category-slug"
				  },
				  "content": "first short content",
				  "createdAt": "2022-04-10 13:32:10.000000011 +0000 UTC"
				},
				{
				  "title": "second title",
				  "slug": "second-title",
				  "category": {
					"name": "second-category",
					"slug": "second-category-slug"
				  },
				  "content": "second short content",
				  "createdAt": "2023-03-05 20:16:30.00000002 +0000 UTC"
				}
			  ],
			  "meta": {
				"offset": 10,
				"limit": 30,
				"total": 1
			  }
			}`,
		},
		{
			name:         "not valid params",
			service:      NewMockArticlesService,
			request:      httptest.NewRequest(http.MethodGet, "/articles?limit=-1&q=test", nil),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"Status":"Bad Request","Messages":["limit must be greater than 0"]}`,
		},
		{
			name:         "empty query",
			service:      NewMockArticlesService,
			request:      httptest.NewRequest(http.MethodGet, "/articles?q=", nil),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"Status":"Bad Request","Messages":["empty query"]}`,
		},
		{
			name: "internal server error",
			service: func(t mockConstructorTestingTNewMockArticlesService) *MockArticlesService {
				m := NewMockArticlesService(t)
				m.On("ArticleListForQuery", mock.Anything, "test-query", int64(10), int64(30)).Return(nil, errors.New("fail get articles"))
				return m
			},
			request:      httptest.NewRequest(http.MethodGet, "/articles?offset=10&limit=30&q=test-query", nil),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"Status":"Internal Server Error","Messages":["internal server error"]}`,
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
