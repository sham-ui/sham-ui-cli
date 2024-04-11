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

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ServeHTTP(t *testing.T) {
	testCases := []struct {
		name                   string
		request                *http.Request
		categoryService        func(t mockConstructorTestingTNewMockCategoryService) *MockCategoryService
		expectedCode           int
		expectedBody           string
		expectedLoggerMessages []testlogger.Message
	}{
		{
			name:    "success",
			request: httptest.NewRequest(http.MethodGet, "/", nil),
			categoryService: func(t mockConstructorTestingTNewMockCategoryService) *MockCategoryService {
				m := NewMockCategoryService(t)
				m.
					On("All", mock.Anything).
					Return([]model.Category{
						{ID: "42", Name: "first category", Slug: "first-category"},
						{ID: "43", Name: "second category", Slug: "second-category"},
					}, nil).
					Once()
				return m
			},
			expectedCode: http.StatusOK,
			expectedBody: `{
			  "categories": [
				{
				  "id": "42",
				  "name": "first category",
				  "slug": "first-category"
				},
				{
				  "id": "43",
				  "name": "second category",
				  "slug": "second-category"
				}
			  ]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:    "failed to get categories",
			request: httptest.NewRequest(http.MethodGet, "/", nil),
			categoryService: func(t mockConstructorTestingTNewMockCategoryService) *MockCategoryService {
				m := NewMockCategoryService(t)
				m.
					On("All", mock.Anything).
					Return(nil, errors.New("test")).
					Once()
				return m
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{
			  "Status": "Internal Server Error",
			  "Messages": ["internal server error"]
			}`,
			expectedLoggerMessages: []testlogger.Message{
				{
					Message: "failed to get categories",
					KeyValues: map[string]interface{}{
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
			Setup(router, test.categoryService(t))
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
