package remove

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
			request: httptest.NewRequest(http.MethodDelete, "/42", nil),
			categoryService: func(t mockConstructorTestingTNewMockCategoryService) *MockCategoryService {
				m := NewMockCategoryService(t)
				m.
					On("Delete", mock.Anything, model.CategoryID("42")).
					Return(nil).
					Once()
				return m
			},
			expectedCode:           http.StatusOK,
			expectedBody:           `{"Status":"Category deleted"}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:    "fail delete category",
			request: httptest.NewRequest(http.MethodDelete, "/42", nil),
			categoryService: func(t mockConstructorTestingTNewMockCategoryService) *MockCategoryService {
				m := NewMockCategoryService(t)
				m.
					On(
						"Delete", mock.Anything, model.CategoryID("42")).
					Return(errors.New("test")).
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
					Message: "fail delete category",
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
