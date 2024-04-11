package detail

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
	article := model.Article{
		ID:          "42",
		Slug:        "article slug",
		Title:       "article title",
		CategoryID:  "12",
		ShortBody:   "Short body",
		Body:        "Full content",
		PublishedAt: time.Date(2024, time.March, 9, 21, 2, 4, 0, time.UTC),
	}
	firstTag := model.Tag{
		ID:   "1",
		Slug: "first-tag",
		Name: "First tag",
	}
	secondTag := model.Tag{
		ID:   "2",
		Slug: "second-tag",
		Name: "Second tag",
	}

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
			request: httptest.NewRequest(http.MethodGet, "/42", nil),
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.EXPECT().
					FindByID(mock.Anything, model.ArticleID("42")).
					Return(article, nil).
					Once()
				m.EXPECT().
					GetTags(mock.Anything, model.ArticleID("42")).
					Return([]model.Tag{firstTag, secondTag}, nil).
					Once()
				return m
			},
			expectedCode: http.StatusOK,
			expectedBody: `{
				"slug": "article slug",
				"title": "article title",
				"category_id": "12",
				"short_body": "Short body",
				"body": "Full content",
				"published_at": "2024-03-09T21:02:04Z",
				"tags": [
					"first-tag",
					"second-tag"
				]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:    "article not found",
			request: httptest.NewRequest(http.MethodGet, "/42", nil),
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.EXPECT().
					FindByID(mock.Anything, model.ArticleID("42")).
					Return(model.Article{}, model.ErrArticleNotFound). //nolint:exhaustruct
					Once()
				return m
			},
			expectedCode:           http.StatusNotFound,
			expectedBody:           `{"Status":"Not Found", "Messages":["Article not found"]}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:    "fail get article",
			request: httptest.NewRequest(http.MethodGet, "/42", nil),
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.EXPECT().
					FindByID(mock.Anything, model.ArticleID("42")).
					Return(model.Article{}, errors.New("test")). //nolint:exhaustruct
					Once()
				return m
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"Status":"Internal Server Error", "Messages":["internal server error"]}`,
			expectedLoggerMessages: []testlogger.Message{
				{
					Message:   "fail find article",
					KeyValues: map[string]any{"error": "test"},
				},
			},
		},
		{
			name:    "fail get article tags",
			request: httptest.NewRequest(http.MethodGet, "/42", nil),
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.EXPECT().
					FindByID(mock.Anything, model.ArticleID("42")).
					Return(article, nil).
					Once()
				m.EXPECT().
					GetTags(mock.Anything, model.ArticleID("42")).
					Return([]model.Tag{}, errors.New("test")).
					Once()
				return m
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"Status":"Internal Server Error", "Messages":["internal server error"]}`,
			expectedLoggerMessages: []testlogger.Message{
				{
					Message: "fail get article tags",
					KeyValues: map[string]any{
						"error": "test",
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(
			test.name, func(t *testing.T) {
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
			},
		)
	}
}
