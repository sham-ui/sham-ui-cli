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
		tagService             func(t mockConstructorTestingTNewMockTagService) *MockTagService
		expectedCode           int
		expectedBody           string
		expectedLoggerMessages []testlogger.Message
	}{
		{
			name:    "success",
			request: httptest.NewRequest(http.MethodGet, "/", nil),
			tagService: func(t mockConstructorTestingTNewMockTagService) *MockTagService {
				m := NewMockTagService(t)
				m.
					On("All", mock.Anything).
					Return([]model.Tag{
						{ID: "42", Name: "first tag", Slug: "first-tag"},
						{ID: "43", Name: "second tag", Slug: "second-tag"},
					}, nil).
					Once()
				return m
			},
			expectedCode: http.StatusOK,
			expectedBody: `{
			  "tags": [
				{
				  "id": "42",
				  "name": "first tag",
				  "slug": "first-tag"
				},
				{
				  "id": "43",
				  "name": "second tag",
				  "slug": "second-tag"
				}
			  ]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:    "failed to get tags",
			request: httptest.NewRequest(http.MethodGet, "/", nil),
			tagService: func(t mockConstructorTestingTNewMockTagService) *MockTagService {
				m := NewMockTagService(t)
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
					Message: "failed to get tags",
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
			Setup(router, test.tagService(t))
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
