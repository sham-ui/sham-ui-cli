package update

import (
	"cms/internal/model"
	"cms/pkg/asserts"
	"cms/pkg/logger"
	"cms/pkg/logger/testlogger"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ServeHTTP(t *testing.T) {
	goodBody := `{
		"title": "new article",
		"category_id": "1",
		"short_body": "new article short body",
		"body": "new article body",
		"published_at": "2022-03-08T05:49:52.643Z",
		"tags": [
			{"name": "existed tag", "slug": "existed-tag"},
			{"name": "new tag"}
		]
	}`
	goodBodyWithoutTags := `{
		"title": "new article",
		"category_id": "1",
		"short_body": "new article short body",
		"body": "new article body",
		"published_at": "2022-03-08T05:49:52.643Z"
	}`
	article := model.Article{
		ID:          "42",
		Slug:        "new-article",
		Title:       "new article",
		CategoryID:  "1",
		ShortBody:   "new article short body",
		Body:        "new article body",
		PublishedAt: time.Date(2022, time.March, 8, 5, 49, 52, 643000000, time.UTC),
	}
	existedTag := model.Tag{ //nolint:exhaustruct
		Slug: "existed-tag",
		Name: "existed tag",
	}
	newTag := model.Tag{ //nolint:exhaustruct
		Slug: "new-tag",
		Name: "new tag",
	}
	tags := []model.Tag{existedTag, newTag}

	testCases := []struct {
		name                   string
		request                *http.Request
		slugifyService         func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService
		articleService         func(t mockConstructorTestingTNewMockArticleService) *MockArticleService
		expectedCode           int
		expectedBody           string
		expectedLoggerMessages []testlogger.Message
	}{
		{
			name:    "success",
			request: httptest.NewRequest(http.MethodPut, "/42", strings.NewReader(goodBody)),
			slugifyService: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.
					On("SlugifyArticle", mock.Anything, "new article").
					Return(article.Slug).
					Once()
				m.
					On("SlugifyTag", mock.Anything, "new tag").
					Return(newTag.Slug).
					Once()
				return m
			},
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("Update", mock.Anything, article, tags).
					Return(nil).
					Once()
				return m
			},
			expectedCode:           http.StatusOK,
			expectedBody:           `{"Status":"Article updated"}`,
			expectedLoggerMessages: []testlogger.Message{},
		},

		{
			name:                   "invalid json",
			request:                httptest.NewRequest(http.MethodPut, "/42", strings.NewReader("{")),
			slugifyService:         NewMockSlugifyService,
			articleService:         NewMockArticleService,
			expectedCode:           http.StatusBadRequest,
			expectedBody:           `{"Status":"Bad Request", "Messages": ["Invalid JSON"]}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:    "article title already exists",
			request: httptest.NewRequest(http.MethodPut, "/42", strings.NewReader(goodBodyWithoutTags)),
			slugifyService: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.
					On("SlugifyArticle", mock.Anything, "new article").
					Return(article.Slug).
					Once()
				return m
			},
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("Update", mock.Anything, article, []model.Tag{}).
					Return(model.ErrArticleTitleAlreadyExists).
					Once()
				return m
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{
				"Status":"Bad Request",
				"Messages":["Title is already in use."]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:    "article slug already exists",
			request: httptest.NewRequest(http.MethodPut, "/42", strings.NewReader(goodBodyWithoutTags)),
			slugifyService: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.
					On("SlugifyArticle", mock.Anything, "new article").
					Return(article.Slug).
					Once()
				return m
			},
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("Update", mock.Anything, article, []model.Tag{}).
					Return(model.ErrArticleSlugAlreadyExists).
					Once()
				return m
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{
				"Status":"Bad Request",
				"Messages":["Slug is already in use."]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:    "category not found",
			request: httptest.NewRequest(http.MethodPut, "/42", strings.NewReader(goodBodyWithoutTags)),
			slugifyService: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.
					On("SlugifyArticle", mock.Anything, "new article").
					Return(article.Slug).
					Once()
				return m
			},
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("Update", mock.Anything, article, []model.Tag{}).
					Return(model.ErrCategoryNotFound).
					Once()
				return m
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{
				"Status":"Bad Request",
				"Messages":["Category not found."]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:    "fail to update",
			request: httptest.NewRequest(http.MethodPut, "/42", strings.NewReader(goodBodyWithoutTags)),
			slugifyService: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.
					On("SlugifyArticle", mock.Anything, "new article").
					Return(article.Slug).
					Once()
				return m
			},
			articleService: func(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
				m := NewMockArticleService(t)
				m.
					On("Update", mock.Anything, article, []model.Tag{}).
					Return(errors.New("test")).
					Once()
				return m
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{
				"Status":"Internal Server Error",
				"Messages":["internal server error"]
			}`,
			expectedLoggerMessages: []testlogger.Message{
				{
					Message: "failed to update article",
					KeyValues: map[string]any{
						"error": "test",
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := mux.NewRouter()
			Setup(router, test.articleService(t), test.slugifyService(t))
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
