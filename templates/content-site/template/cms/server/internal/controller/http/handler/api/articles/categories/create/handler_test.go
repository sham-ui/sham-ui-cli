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
		categoryService        func(t mockConstructorTestingTNewMockCategoryService) *MockCategoryService
		expectedCode           int
		expectedBody           string
		expectedLoggerMessages []testlogger.Message
	}{
		{
			name: "success",
			request: httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{"name": "new category"}`),
			),
			slugifyService: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.
					On("SlugifyCategory", mock.Anything, "new category").
					Return(model.CategorySlug("new-category")).
					Once()
				return m
			},
			categoryService: func(t mockConstructorTestingTNewMockCategoryService) *MockCategoryService {
				m := NewMockCategoryService(t)
				m.
					On("Create", mock.Anything, model.Category{ //nolint:exhaustruct
						Slug: "new-category",
						Name: "new category",
					}).
					Return(nil).
					Once()
				return m
			},
			expectedCode:           http.StatusCreated,
			expectedBody:           `{"Status":"Category created"}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "invalid json",
			request: httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{"name": "new category`),
			),
			slugifyService:         NewMockSlugifyService,
			categoryService:        NewMockCategoryService,
			expectedCode:           http.StatusBadRequest,
			expectedBody:           `{"Status":"Bad Request","Messages":["Invalid JSON"]}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "slug is already in use",
			request: httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{"name": "new category"}`),
			),
			slugifyService: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.
					On("SlugifyCategory", mock.Anything, "new category").
					Return(model.CategorySlug("new-category")).
					Once()
				return m
			},
			categoryService: func(t mockConstructorTestingTNewMockCategoryService) *MockCategoryService {
				m := NewMockCategoryService(t)
				m.
					On("Create", mock.Anything, model.Category{ //nolint:exhaustruct
						Slug: "new-category",
						Name: "new category",
					}).
					Return(model.ErrCategorySlugAlreadyExists).
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
				strings.NewReader(`{"name": "new category"}`),
			),
			slugifyService: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.
					On("SlugifyCategory", mock.Anything, "new category").
					Return(model.CategorySlug("new-category")).
					Once()
				return m
			},
			categoryService: func(t mockConstructorTestingTNewMockCategoryService) *MockCategoryService {
				m := NewMockCategoryService(t)
				m.
					On("Create", mock.Anything, model.Category{ //nolint:exhaustruct
						Slug: "new-category",
						Name: "new category",
					}).
					Return(model.ErrCategoryNameAlreadyExists).
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
			name: "failed to create category",
			request: httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{"name": "new category"}`),
			),
			slugifyService: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.
					On("SlugifyCategory", mock.Anything, "new category").
					Return(model.CategorySlug("new-category")).
					Once()
				return m
			},
			categoryService: func(t mockConstructorTestingTNewMockCategoryService) *MockCategoryService {
				m := NewMockCategoryService(t)
				m.
					On("Create", mock.Anything, model.Category{ //nolint:exhaustruct
						Slug: "new-category",
						Name: "new category",
					}).
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
					Message:   "failed to create category",
					KeyValues: map[string]any{"error": "test"},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := mux.NewRouter()
			Setup(router, test.categoryService(t), test.slugifyService(t))
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
