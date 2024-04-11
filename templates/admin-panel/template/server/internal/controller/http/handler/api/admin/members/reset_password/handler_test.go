package reset_password

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"{{ shortName }}/internal/model"
	"{{ shortName }}/pkg/asserts"
	"{{ shortName }}/pkg/logger"
	"{{ shortName }}/pkg/logger/testlogger"
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
		passwordService        func(t mockConstructorTestingTNewMockPasswordService) *MockPasswordService
		memberService          func(t mockConstructorTestingTNewMockMemberService) *MockMemberService
		expectedCode           int
		expectedBody           string
		expectedLoggerMessages []testlogger.Message
	}{
		{
			name: "success",
			request: httptest.NewRequest(
				http.MethodPut,
				"/42/password",
				strings.NewReader(`{
					"pass1": "password1", 
					"pass2": "password1" 
				}`),
			),
			passwordService: func(t mockConstructorTestingTNewMockPasswordService) *MockPasswordService {
				m := NewMockPasswordService(t)
				m.
					On("Hash", mock.Anything, "password1").
					Return(model.MemberHashedPassword("hashed-pass"), nil).
					Once()
				return m
			},
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("UpdatePassword", mock.Anything, model.MemberID("42"), model.MemberHashedPassword("hashed-pass")).
					Return(nil).
					Once()
				return m
			},
			expectedCode:           http.StatusOK,
			expectedBody:           `{"Status":"Password updated"}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "json decode error",
			request: httptest.NewRequest(
				http.MethodPut,
				"/42/password",
				strings.NewReader(""),
			),
			passwordService:        NewMockPasswordService,
			memberService:          NewMockMemberService,
			expectedCode:           http.StatusBadRequest,
			expectedBody:           `{"Status": "Bad Request", "Messages": ["Invalid JSON"]}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "pass1 mismatch",
			request: httptest.NewRequest(
				http.MethodPut,
				"/42/password",
				strings.NewReader(`{"pass2": "new-pass"}`),
			),
			passwordService: NewMockPasswordService,
			memberService:   NewMockMemberService,
			expectedCode:    http.StatusBadRequest,
			expectedBody: `{
				"Status": "Bad Request",
				"Messages": ["Password must have more than 0 characters.", "Passwords don't match."]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "pass2 mismatch",
			request: httptest.NewRequest(
				http.MethodPut,
				"/42/password",
				strings.NewReader(`{"pass1": "new-pass"}`),
			),
			passwordService: NewMockPasswordService,
			memberService:   NewMockMemberService,
			expectedCode:    http.StatusBadRequest,
			expectedBody: `{
				"Status": "Bad Request",
				"Messages": ["Password must have more than 0 characters.", "Passwords don't match."]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "passwords don't match",
			request: httptest.NewRequest(
				http.MethodPut,
				"/42/password",
				strings.NewReader(`{"pass1": "new-pass", "pass2": "new-pass-2"}`),
			),
			passwordService: NewMockPasswordService,
			memberService:   NewMockMemberService,
			expectedCode:    http.StatusBadRequest,
			expectedBody: `{
				"Status": "Bad Request",
				"Messages": ["Passwords don't match."]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "fail hash password",
			request: httptest.NewRequest(
				http.MethodPut,
				"/42/password",
				strings.NewReader(`{"pass1": "new-pass", "pass2": "new-pass"}`),
			),
			passwordService: func(t mockConstructorTestingTNewMockPasswordService) *MockPasswordService {
				m := NewMockPasswordService(t)
				m.
					On("Hash", mock.Anything, "new-pass").
					Return(model.MemberHashedPassword(""), errors.New("test")).
					Once()
				return m
			},
			memberService: NewMockMemberService,
			expectedCode:  http.StatusInternalServerError,
			expectedBody: `{
				"Status": "Internal Server Error",
				"Messages": ["internal server error"]
			}`,
			expectedLoggerMessages: []testlogger.Message{
				{
					Message: "fail hash password",
					KeyValues: map[string]any{
						"error": "test",
					},
				},
			},
		},
		{
			name: "fail update password",
			request: httptest.NewRequest(
				http.MethodPut,
				"/42/password",
				strings.NewReader(`{"pass1": "new-pass", "pass2": "new-pass"}`),
			),
			passwordService: func(t mockConstructorTestingTNewMockPasswordService) *MockPasswordService {
				m := NewMockPasswordService(t)
				m.
					On("Hash", mock.Anything, "new-pass").
					Return(model.MemberHashedPassword("hashed-pass"), nil).
					Once()
				return m
			},
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On(
						"UpdatePassword",
						mock.Anything,
						model.MemberID("42"),
						model.MemberHashedPassword("hashed-pass"),
					).
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
					Message: "fail update member password",
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
			Setup(router, test.memberService(t), test.passwordService(t))
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
