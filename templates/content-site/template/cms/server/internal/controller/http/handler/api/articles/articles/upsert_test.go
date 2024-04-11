package articles

import (
	"bytes"
	"cms/internal/model"
	"cms/pkg/asserts"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestExtractAndValidateData(t *testing.T) {
	testCases := []struct {
		name              string
		req               *http.Request
		slugger           func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService
		expectedArticle   *model.Article
		expectedTags      []model.Tag
		expectedDataValid bool
	}{
		{
			name: "success",
			req: httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`{
				"title": "new article",
				"category_id": "1",
				"short_body": "new article short body",
				"body": "new article body",
				"published_at": "2022-03-08T05:49:52.643Z",
				"tags": [
					{"name": "existed tag", "slug": "existed-tag"},
					{"name": "new tag"}
				]
			}`))),
			slugger: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.EXPECT().
					SlugifyArticle(mock.Anything, "new article").
					Return("new-article").
					Once()
				m.EXPECT().
					SlugifyTag(mock.Anything, "new tag").
					Return("new-tag").
					Once()
				return m
			},
			expectedArticle: &model.Article{ //nolint:exhaustruct
				Slug:        "new-article",
				Title:       "new article",
				CategoryID:  "1",
				ShortBody:   "new article short body",
				Body:        "new article body",
				PublishedAt: time.Date(2022, time.March, 8, 5, 49, 52, 643000000, time.UTC),
			},
			expectedTags: []model.Tag{
				{
					Slug: "existed-tag",
					Name: "existed tag",
				},
				{
					Slug: "new-tag",
					Name: "new tag",
				},
			},
			expectedDataValid: true,
		},
		{
			name:              "fail parse json",
			req:               httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(``))),
			slugger:           NewMockSlugifyService,
			expectedDataValid: false,
		},
		{
			name:              "empty title",
			req:               httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`{}`))),
			slugger:           NewMockSlugifyService,
			expectedDataValid: false,
		},
		{
			name: "empty category_id",
			req: httptest.NewRequest(
				http.MethodPost,
				"/",
				bytes.NewBuffer([]byte(`{"title": "test"'}`)),
			),
			slugger:           NewMockSlugifyService,
			expectedDataValid: false,
		},
		{
			name: "trim data",
			req: httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`{
				"title": "    new article    ",
				"category_id": "1",
				"short_body": "    new article short body    ",
				"body": "    new article body    ",
				"published_at": "2022-03-08T05:49:52.643Z",
				"tags": [
					{"name": "existed tag", "slug": "existed-tag"},
					{"name": "    new tag    "}
				]
			}`))),
			slugger: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.EXPECT().
					SlugifyArticle(mock.Anything, "new article").
					Return("new-article").
					Once()
				m.EXPECT().
					SlugifyTag(mock.Anything, "new tag").
					Return("new-tag").
					Once()
				return m
			},
			expectedArticle: &model.Article{ //nolint:exhaustruct
				Slug:        "new-article",
				Title:       "new article",
				CategoryID:  "1",
				ShortBody:   "new article short body",
				Body:        "new article body",
				PublishedAt: time.Date(2022, time.March, 8, 5, 49, 52, 643000000, time.UTC),
			},
			expectedTags: []model.Tag{
				{
					Slug: "existed-tag",
					Name: "existed tag",
				},
				{
					Slug: "new-tag",
					Name: "new tag",
				},
			},
			expectedDataValid: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			rw := httptest.NewRecorder()

			// Action
			article, tags, valid := ExtractAndValidateData(test.slugger(t), rw, test.req)

			// Assert
			asserts.Equals(t, test.expectedArticle, article)
			asserts.Equals(t, test.expectedTags, tags)
			asserts.Equals(t, test.expectedDataValid, valid)
		})
	}
}

func TestHandleError(t *testing.T) {
	testCases := []struct {
		name               string
		err                error
		expectedHandled    bool
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "title already exists",
			err:                model.ErrArticleTitleAlreadyExists,
			expectedHandled:    true,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"Status": "Bad Request", "Messages": ["Title is already in use."]}`,
		},
		{
			name:               "slug already exists",
			err:                model.ErrArticleSlugAlreadyExists,
			expectedHandled:    true,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"Status": "Bad Request", "Messages": ["Slug is already in use."]}`,
		},
		{
			name:               "category not found",
			err:                model.ErrCategoryNotFound,
			expectedHandled:    true,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"Status": "Bad Request", "Messages": ["Category not found."]}`,
		},
		{
			name:               "other error",
			err:                errors.New("test"),
			expectedHandled:    false,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "err == nil",
			err:                nil,
			expectedHandled:    false,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			rw := httptest.NewRecorder()

			// Action
			handled := HandleError(test.err, rw, r)

			// Assert
			asserts.Equals(t, test.expectedHandled, handled)
			asserts.Equals(t, test.expectedStatusCode, rw.Code)
			if handled {
				asserts.JSONEquals(t, test.expectedBody, rw.Body.String())
			}
		})
	}
}
