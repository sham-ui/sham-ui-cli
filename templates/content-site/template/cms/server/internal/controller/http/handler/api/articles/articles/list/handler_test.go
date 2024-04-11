package list

import (
	"cms/internal/model"
	"cms/pkg/asserts"
	"cms/pkg/logger"
	"cms/pkg/logger/testlogger"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ServeHTTP(t *testing.T) {
	testCases := []struct {
		name                   string
		request                *http.Request
		articleService         func(t mockConstructorTestingTNewMockArticleService) *MockArticleService
		expectedCode           int
		expectedBody           string
		expectedLoggerMessages []testlogger.Message
	}{
		{
			name:    "success",
			request: httptest.NewRequest(http.MethodGet, "/?offset=10&limit=30", nil),
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("FindShortInfo", mock.Anything, int64(10), int64(30)).
					Return([]model.ArticleShortInfo{
						{
							ID:          "42",
							Slug:        "first-article",
							Title:       "first article",
							CategoryID:  "1",
							PublishedAt: time.Date(2020, time.January, 1, 2, 32, 4, 0, time.UTC),
						},
						{
							ID:          "43",
							Slug:        "second-article",
							Title:       "second article",
							CategoryID:  "2",
							PublishedAt: time.Date(2023, time.March, 9, 21, 2, 4, 0, time.UTC),
						},
					}, nil).
					Once()
				m.
					On("Total", mock.Anything).
					Return(2, nil).
					Once()
				return m
			},
			expectedCode: http.StatusOK,
			expectedBody: `{
			  "articles": [
				{
				  "id": "42",
				  "slug": "first-article",	
				  "title": "first article",
				  "category_id": "1",
				  "published_at": "2020-01-01T02:32:04Z"
				},
				{
				  "id": "43",
				  "slug": "second-article",
				  "title": "second article",
				  "category_id": "2",
				  "published_at": "2023-03-09T21:02:04Z"
				}
			  ],
			  "meta": {
				"offset": 10,
				"limit": 30,
				"total": 2
			  }
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:    "empty list",
			request: httptest.NewRequest(http.MethodGet, "/?offset=10&limit=30", nil),
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("FindShortInfo", mock.Anything, int64(10), int64(30)).
					Return([]model.ArticleShortInfo{}, nil).
					Once()
				m.
					On("Total", mock.Anything).
					Return(0, nil).
					Once()
				return m
			},
			expectedCode: http.StatusOK,
			expectedBody: `{
			  "articles": [],
			  "meta": {
				"offset": 10,
				"limit": 30,
				"total": 0
			  }
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:           "limit=0",
			request:        httptest.NewRequest(http.MethodGet, "/?offset=abc&limit=0", nil),
			articleService: NewMockArticleService,
			expectedCode:   http.StatusBadRequest,
			expectedBody: `{
			  "Status": "Bad Request",
			  "Messages": [
				"limit must be greater than 0"
			  ]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:    "failed to find articles",
			request: httptest.NewRequest(http.MethodGet, "/?offset=10&limit=30", nil),
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("FindShortInfo", mock.Anything, int64(10), int64(30)).
					Return(nil, errors.New("test")).
					Once()
				return m
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{
			  "Status": "Internal Server Error",
			  "Messages": [
				"internal server error"
			  ]
			}`,
			expectedLoggerMessages: []testlogger.Message{
				{
					Message:   "failed to find articles",
					KeyValues: map[string]interface{}{"error": "test"},
				},
			},
		},
		{
			name:    "failed to count articles",
			request: httptest.NewRequest(http.MethodGet, "/?offset=10&limit=30", nil),
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("FindShortInfo", mock.Anything, int64(10), int64(30)).
					Return([]model.ArticleShortInfo{}, nil).
					Once()
				m.
					On("Total", mock.Anything).
					Return(0, errors.New("test")).
					Once()
				return m
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{
			  "Status": "Internal Server Error",
			  "Messages": [
				"internal server error"
			  ]
			}`,
			expectedLoggerMessages: []testlogger.Message{
				{
					Message:   "failed to count articles",
					KeyValues: map[string]interface{}{"error": "test"},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := mux.NewRouter()
			Setup(router, test.articleService(t))
			log := testlogger.NewLogger()
			ctx := logger.Save(test.request.Context(), logr.New(log))
			req := test.request.WithContext(ctx)
			resp := httptest.NewRecorder()

			// Action
			router.ServeHTTP(resp, req)

			// Assert
			asserts.Equals(t, test.expectedCode, resp.Code, "code")
			asserts.JSONEquals(t, test.expectedBody, resp.Body.String(), "body")
			asserts.Equals(t, test.expectedLoggerMessages, log.Messages, "logger")
		})
	}
}
