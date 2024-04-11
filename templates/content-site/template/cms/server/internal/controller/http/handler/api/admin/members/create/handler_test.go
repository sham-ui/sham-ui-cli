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
		passwordService        func(t mockConstructorTestingTNewMockPasswordService) *MockPasswordService
		memberService          func(t mockConstructorTestingTNewMockMemberService) *MockMemberService
		expectedCode           int
		expectedBody           string
		expectedLoggerMessages []testlogger.Message
	}{
		{
			name: "success",
			request: httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{
					"email": "test@example.com", 
					"name": "tester", 
					"password": "test", 
					"is_superuser": true
				}`),
			),
			passwordService: func(t mockConstructorTestingTNewMockPasswordService) *MockPasswordService {
				m := NewMockPasswordService(t)
				m.
					On("Hash", mock.Anything, mock.Anything).
					Return(model.MemberHashedPassword("hashed"), nil).
					Once()
				return m
			},
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("Create", mock.Anything, model.MemberWithPassword{
						Member: model.Member{ //nolint:exhaustruct
							Email:       "test@example.com",
							Name:        "tester",
							IsSuperuser: true,
						},
						HashedPassword: model.MemberHashedPassword("hashed"),
					}).
					Return(nil).
					Once()
				return m
			},
			expectedCode:           http.StatusOK,
			expectedBody:           `{"Status":"Member created"}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:                   "json decode error",
			request:                httptest.NewRequest(http.MethodPost, "/", strings.NewReader("")),
			passwordService:        NewMockPasswordService,
			memberService:          NewMockMemberService,
			expectedCode:           http.StatusBadRequest,
			expectedBody:           `{"Status": "Bad Request", "Messages": ["Invalid JSON"]}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "validation error",
			request: httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{}`),
			),
			passwordService: NewMockPasswordService,
			memberService:   NewMockMemberService,
			expectedCode:    http.StatusBadRequest,
			expectedBody: `{
				"Status": "Bad Request",
				"Messages": [
					"Name must not be empty.",
					"Email must not be empty.",
					"Password must not be empty."
				]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "hash password error",
			request: httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{
					"email": "test@example.com", 
					"name": "tester", 
					"password": "password", 
					"is_superuser": true
				}`),
			),
			passwordService: func(t mockConstructorTestingTNewMockPasswordService) *MockPasswordService {
				m := NewMockPasswordService(t)
				m.
					On("Hash", mock.Anything, "password").
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
					Message:   "failed to hash password",
					KeyValues: map[string]any{"error": "test"},
				},
			},
		},
		{
			name: "email already exists",
			request: httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{
					"email": "test@example.com", 
					"name": "tester", 
					"password": "password", 
					"is_superuser": true
				}`),
			),
			passwordService: func(t mockConstructorTestingTNewMockPasswordService) *MockPasswordService {
				m := NewMockPasswordService(t)
				m.
					On("Hash", mock.Anything, "password").
					Return(model.MemberHashedPassword("hashed"), nil).
					Once()
				return m
			},
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("Create", mock.Anything, model.MemberWithPassword{
						Member: model.Member{ //nolint:exhaustruct
							Email:       "test@example.com",
							Name:        "tester",
							IsSuperuser: true,
						},
						HashedPassword: model.MemberHashedPassword("hashed"),
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
			name: "create error",
			request: httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{
					"email": "test@example.com", 
					"name": "tester", 
					"password": "password", 
					"is_superuser": true
				}`),
			),
			passwordService: func(t mockConstructorTestingTNewMockPasswordService) *MockPasswordService {
				m := NewMockPasswordService(t)
				m.
					On("Hash", mock.Anything, "password").
					Return(model.MemberHashedPassword("hashed"), nil).
					Once()
				return m
			},
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("Create", mock.Anything, model.MemberWithPassword{
						Member: model.Member{ //nolint:exhaustruct
							Email:       "test@example.com",
							Name:        "tester",
							IsSuperuser: true,
						},
						HashedPassword: model.MemberHashedPassword("hashed"),
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
					Message:   "failed to create member",
					KeyValues: map[string]any{"error": "test"},
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
