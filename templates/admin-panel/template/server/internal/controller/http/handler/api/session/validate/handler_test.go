package validate

import (
	"context"
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
)

func TestHandler(t *testing.T) {
	testCases := []struct {
		name                   string
		request                *http.Request
		expectedCode           int
		expectedBody           string
		expectedLoggerMessages []testlogger.Message
	}{
		{
			name: "success",
			request: httptest.NewRequest(http.MethodGet, "/validsession", nil).
				WithContext(request.SaveSessionToContext(context.Background(), &model.Session{
					MemberID:    "42",
					Email:       "test@example.com",
					Name:        "tester",
					IsSuperuser: true,
				})),
			expectedCode:           http.StatusOK,
			expectedBody:           "{\"Name\":\"tester\",\"Email\":\"test@example.com\",\"IsSuperuser\":true}\n",
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:                   "not authenticated",
			request:                httptest.NewRequest(http.MethodGet, "/validsession", nil).WithContext(context.Background()),
			expectedCode:           http.StatusUnauthorized,
			expectedBody:           "{\"Status\":\"Unauthorized\",\"Messages\":[\"not authenticated\"]}\n",
			expectedLoggerMessages: []testlogger.Message{},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := mux.NewRouter()
			Setup(router)
			log := testlogger.NewLogger()
			ctx := logger.Save(test.request.Context(), logr.New(log))
			req := test.request.WithContext(ctx)
			resp := httptest.NewRecorder()

			// Action
			router.ServeHTTP(resp, req)

			// Assert
			asserts.JSONEquals(t, test.expectedBody, resp.Body.String(), "body")
			asserts.Equals(t, test.expectedCode, resp.Code, "code")
			asserts.Equals(t, test.expectedLoggerMessages, log.Messages, "logger")
		})
	}
}
