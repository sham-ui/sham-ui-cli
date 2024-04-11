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

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ServeHTTP(t *testing.T) {
	testCases := []struct {
		name                   string
		request                *http.Request
		memberService          func(t mockConstructorTestingTNewMockMemberService) *MockMemberService
		expectedCode           int
		expectedBody           string
		expectedLoggerMessages []testlogger.Message
	}{
		{
			name: "success",
			request: httptest.NewRequest(
				http.MethodPut,
				"/42",
				strings.NewReader(`{
					"email": "test@example.com", 
					"name": "tester", 
					"password": "test", 
					"is_superuser": true
				}`),
			),
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("Update", mock.Anything, model.Member{
						ID:          "42",
						Email:       "test@example.com",
						Name:        "tester",
						IsSuperuser: true,
					}).
					Return(nil).
					Once()
				return m
			},
			expectedCode:           http.StatusOK,
			expectedBody:           `{"Status":"Member updated"}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:                   "json decode error",
			request:                httptest.NewRequest(http.MethodPut, "/42", strings.NewReader("")),
			memberService:          NewMockMemberService,
			expectedCode:           http.StatusBadRequest,
			expectedBody:           `{"Status": "Bad Request", "Messages": ["Invalid JSON"]}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "validation error",
			request: httptest.NewRequest(
				http.MethodPut,
				"/42",
				strings.NewReader(`{}`),
			),
			memberService: NewMockMemberService,
			expectedCode:  http.StatusBadRequest,
			expectedBody: `{
				"Status": "Bad Request",
				"Messages": [
					"Name must not be empty.",
					"Email must not be empty."
				]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "email already exists",
			request: httptest.NewRequest(
				http.MethodPut,
				"/42",
				strings.NewReader(`{
					"email": "test@example.com", 
					"name": "tester", 
					"password": "password", 
					"is_superuser": true
				}`),
			),
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("Update", mock.Anything, model.Member{
						ID:          "42",
						Email:       "test@example.com",
						Name:        "tester",
						IsSuperuser: true,
					}).
					Return(model.ErrMemberEmailAlreadyExists).
					Once()
				return m
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{
				"Status": "Bad Request",
				"Messages": ["Email is already in use."]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "update error",
			request: httptest.NewRequest(
				http.MethodPut,
				"/42",
				strings.NewReader(`{
					"email": "test@example.com", 
					"name": "tester", 
					"password": "password", 
					"is_superuser": true
				}`),
			),
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("Update", mock.Anything, model.Member{
						ID:          "42",
						Email:       "test@example.com",
						Name:        "tester",
						IsSuperuser: true,
					}).
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
					Message:   "failed to update member",
					KeyValues: map[string]any{"error": "test"},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := mux.NewRouter()
			Setup(router, test.memberService(t))
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
