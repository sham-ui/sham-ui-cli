package create

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

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ServeHTTP(t *testing.T) {
	testCases := []struct {
		name                   string
		request                *http.Request
		slugifyService         func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService
		tagService             func(t mockConstructorTestingTNewMockTagService) *MockTagService
		expectedCode           int
		expectedBody           string
		expectedLoggerMessages []testlogger.Message
	}{
		{
			name: "success",
			request: httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{"name": "new tag"}`),
			),
			slugifyService: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.
					On("SlugifyTag", mock.Anything, "new tag").
					Return(model.TagSlug("new-tag")).
					Once()
				return m
			},
			tagService: func(t mockConstructorTestingTNewMockTagService) *MockTagService {
				m := NewMockTagService(t)
				m.
					On("Create", mock.Anything, model.Tag{ //nolint:exhaustruct
						Slug: "new-tag",
						Name: "new tag",
					}).
					Return(model.TagID("1"), nil).
					Once()
				return m
			},
			expectedCode:           http.StatusCreated,
			expectedBody:           `{"Status":"Tag created"}`,
			expectedLoggerMessages: []testlogger.Message{},
		},

		{
			name: "invalid json",
			request: httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{"name": "new tag`),
			),
			slugifyService:         NewMockSlugifyService,
			tagService:             NewMockTagService,
			expectedCode:           http.StatusBadRequest,
			expectedBody:           `{"Status":"Bad Request","Messages":["Invalid JSON"]}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "slug is already in use",
			request: httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{"name": "new tag"}`),
			),
			slugifyService: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.
					On("SlugifyTag", mock.Anything, "new tag").
					Return(model.TagSlug("new-tag")).
					Once()
				return m
			},
			tagService: func(t mockConstructorTestingTNewMockTagService) *MockTagService {
				m := NewMockTagService(t)
				m.
					On("Create", mock.Anything, model.Tag{ //nolint:exhaustruct
						Slug: "new-tag",
						Name: "new tag",
					}).
					Return(model.TagID(""), model.ErrTagSlugAlreadyExists).
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
			name: "name is already in use",
			request: httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{"name": "new tag"}`),
			),
			slugifyService: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.
					On("SlugifyTag", mock.Anything, "new tag").
					Return(model.TagSlug("new-tag")).
					Once()
				return m
			},
			tagService: func(t mockConstructorTestingTNewMockTagService) *MockTagService {
				m := NewMockTagService(t)
				m.
					On("Create", mock.Anything, model.Tag{ //nolint:exhaustruct
						Slug: "new-tag",
						Name: "new tag",
					}).
					Return(model.TagID(""), model.ErrTagNameAlreadyExists).
					Once()
				return m
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{
				"Status":"Bad Request",
				"Messages":["Name is already in use."]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "failed to create tag",
			request: httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{"name": "new tag"}`),
			),
			slugifyService: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.
					On("SlugifyTag", mock.Anything, "new tag").
					Return(model.TagSlug("new-tag")).
					Once()
				return m
			},
			tagService: func(t mockConstructorTestingTNewMockTagService) *MockTagService {
				m := NewMockTagService(t)
				m.
					On("Create", mock.Anything, model.Tag{ //nolint:exhaustruct
						Slug: "new-tag",
						Name: "new tag",
					}).
					Return(model.TagID(""), errors.New("test")).
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
					Message:   "failed to create tag",
					KeyValues: map[string]any{"error": "test"},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := mux.NewRouter()
			Setup(router, test.tagService(t), test.slugifyService(t))
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
