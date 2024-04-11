package logout

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"{{ shortName }}/internal/controller/http/request"
	"{{ shortName }}/internal/model"
	"{{ shortName }}/pkg/asserts"
	"{{ shortName }}/pkg/logger"
	"{{ shortName }}/pkg/logger/testlogger"
	"testing"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ServeHTTP(t *testing.T) {
	authenticatedCtx := request.SaveSessionToContext(context.Background(), &model.Session{
		MemberID:    "42",
		Email:       "test@example.com",
		Name:        "tester",
		IsSuperuser: true,
	})

	nonAuthenticatedCtx := context.Background()

	testCases := []struct {
		name                   string
		service                func(t mockConstructorTestingTNewMockService) *MockService
		request                *http.Request
		expectedCode           int
		expectedBody           string
		expectedLoggerMessages []testlogger.Message
	}{
		{
			name: "success",
			service: func(t mockConstructorTestingTNewMockService) *MockService {
				m := NewMockService(t)
				m.On("Delete", mock.Anything, mock.Anything).Return(nil)
				return m
			},
			request: httptest.NewRequest(http.MethodPost, "/logout", nil).
				WithContext(authenticatedCtx),
			expectedCode:           http.StatusOK,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:    "not authenticated",
			service: NewMockService,
			request: httptest.NewRequest(http.MethodPost, "/logout", nil).
				WithContext(nonAuthenticatedCtx),
			expectedCode: http.StatusForbidden,
			expectedBody: "{\"Status\":\"Forbidden\",\"Messages\":[\"not authenticated\"]}\n",
			expectedLoggerMessages: []testlogger.Message{
				{
					Message: "session not authenticated",
					KeyValues: map[string]any{
						"error": "session not authenticated",
					},
				},
			},
		},
		{
			name: "internal server error",
			service: func(t mockConstructorTestingTNewMockService) *MockService {
				m := NewMockService(t)
				m.On("Delete", mock.Anything, mock.Anything).Return(errors.New("test"))
				return m
			},
			request: httptest.NewRequest(http.MethodPost, "/logout", nil).
				WithContext(authenticatedCtx),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "{\"Status\":\"Internal Server Error\",\"Messages\":[\"internal server error\"]}\n",
			expectedLoggerMessages: []testlogger.Message{
				{
					Message: "failed to delete session",
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
			srv := test.service(t)
			Setup(router, srv)
			log := testlogger.NewLogger()
			ctx := logger.Save(test.request.Context(), logr.New(log))
			req := test.request.WithContext(ctx)
			resp := httptest.NewRecorder()

			// Action
			router.ServeHTTP(resp, req)

			// Assert
			asserts.Equals(t, test.expectedBody, resp.Body.String(), "body")
			asserts.Equals(t, test.expectedCode, resp.Code, "code")
			asserts.Equals(t, test.expectedLoggerMessages, log.Messages, "logger")
		})
	}
}
