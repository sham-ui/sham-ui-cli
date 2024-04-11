package login

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
	"strings"
	"testing"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ServeHTTP(t *testing.T) {
	memberWithPassword := model.MemberWithPassword{
		Member: model.Member{
			ID:          "42",
			Email:       "test@example.com",
			Name:        "tester",
			IsSuperuser: true,
		},
		// hash from "passw0rd"
		HashedPassword: model.MemberHashedPassword("$2a$14$nsbhfJ5GcjWFOD9S57BXFuxNt1kTmjqaCq7BDRAf8vFHq3.Qgbb8y"),
	}

	testCases := []struct {
		name                   string
		request                *http.Request
		sessionService         func(t mockConstructorTestingTNewMockSessionService) *MockSessionService
		memberService          func(t mockConstructorTestingTNewMockMemberService) *MockMemberService
		passwordService        func(t mockConstructorTestingTNewMockPasswordService) *MockPasswordService
		expectedCode           int
		expectedBody           string
		expectedLoggerMessages []testlogger.Message
	}{
		{
			name: "success",
			request: httptest.NewRequest(
				http.MethodPost,
				"/login",
				strings.NewReader(`{"Email": "test@example.com", "Password": "passw0rd"}`),
			),
			sessionService: func(t mockConstructorTestingTNewMockSessionService) *MockSessionService {
				m := NewMockSessionService(t)
				m.
					On("Create", mock.Anything, mock.Anything, &memberWithPassword.Member).
					Return(nil).
					Once()
				return m
			},
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(&memberWithPassword, nil).Once()
				return m
			},
			passwordService: func(t mockConstructorTestingTNewMockPasswordService) *MockPasswordService {
				m := NewMockPasswordService(t)
				m.
					On("Compare", mock.Anything, memberWithPassword.HashedPassword, "passw0rd").
					Return(nil).
					Once()
				return m
			},
			expectedCode: http.StatusOK,
			expectedBody: `{
				"Status": "OK",
				"Name": "tester",
				"Email": "test@example.com",
				"IsSuperuser": true
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "password mismatch",
			request: httptest.NewRequest(
				http.MethodPost,
				"/login",
				strings.NewReader(`{"Email": "test@example.com", "Password": "invalidpass"}`),
			),
			sessionService: NewMockSessionService,
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(&memberWithPassword, nil).Once()
				return m
			},
			passwordService: func(t mockConstructorTestingTNewMockPasswordService) *MockPasswordService {
				m := NewMockPasswordService(t)
				m.
					On("Compare", mock.Anything, memberWithPassword.HashedPassword, "invalidpass").
					Return(errors.New("incorrect password")).
					Once()
				return m
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{
				"Status": "Bad Request",
				"Messages": ["Incorrect username or password"]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "member not found",
			request: httptest.NewRequest(
				http.MethodPost,
				"/login",
				strings.NewReader(`{"Email": "test@example.com", "Password": "passw0rd"}`),
			),
			passwordService: NewMockPasswordService,
			sessionService:  NewMockSessionService,
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("GetByEmail", mock.Anything, "test@example.com").
					Return(nil, model.ErrMemberNotFound).
					Once()
				return m
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{
				"Status": "Bad Request",
				"Messages": ["Member not found"]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "fail get member",
			request: httptest.NewRequest(
				http.MethodPost,
				"/login",
				strings.NewReader(`{"Email": "test@example.com", "Password": "passw0rd"}`),
			),
			passwordService: NewMockPasswordService,
			sessionService:  NewMockSessionService,
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("GetByEmail", mock.Anything, "test@example.com").
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
					Message: "fail get member",
					KeyValues: map[string]any{
						"error": "test",
					},
				},
			},
		},
		{
			name: "fail create session",
			request: httptest.NewRequest(
				http.MethodPost,
				"/login",
				strings.NewReader(`{"Email": "test@example.com", "Password": "passw0rd"}`),
			),
			passwordService: func(t mockConstructorTestingTNewMockPasswordService) *MockPasswordService {
				m := NewMockPasswordService(t)
				m.
					On("Compare", mock.Anything, memberWithPassword.HashedPassword, "passw0rd").
					Return(nil).
					Once()
				return m
			},
			sessionService: func(t mockConstructorTestingTNewMockSessionService) *MockSessionService {
				m := NewMockSessionService(t)
				m.
					On("Create", mock.Anything, mock.Anything, &memberWithPassword.Member).
					Return(errors.New("test")).
					Once()
				return m
			},
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("GetByEmail", mock.Anything, "test@example.com").
					Return(&memberWithPassword, nil).
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
					Message: "fail create session",
					KeyValues: map[string]any{
						"error": "test",
					},
				},
			},
		},
		{
			name: "already logined",
			request: httptest.NewRequest(
				http.MethodPost,
				"/login",
				strings.NewReader(`{"Email": "test@example.com", "Password": "passw0rd"}`),
			).WithContext(
				request.SaveSessionToContext(context.Background(), &model.Session{
					MemberID:    memberWithPassword.ID,
					Email:       memberWithPassword.Email,
					Name:        memberWithPassword.Name,
					IsSuperuser: memberWithPassword.IsSuperuser,
				})),
			passwordService: NewMockPasswordService,
			sessionService:  NewMockSessionService,
			memberService:   NewMockMemberService,
			expectedCode:    http.StatusBadRequest,
			expectedBody: `{
				"Status": "Bad Request",
				"Messages": ["Already logged in"]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := mux.NewRouter()
			Setup(
				router,
				test.sessionService(t),
				test.memberService(t),
				test.passwordService(t),
				"X-CSRF-Token",
			)
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
