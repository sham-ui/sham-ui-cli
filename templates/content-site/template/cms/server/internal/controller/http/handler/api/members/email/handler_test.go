package email

import (
	"cms/internal/controller/http/request"
	"cms/internal/model"
	"cms/pkg/asserts"
	"cms/pkg/logger"
	"cms/pkg/logger/testlogger"
	"context"
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
	contextWithSession := request.SaveSessionToContext(context.Background(), &model.Session{
		MemberID:    "42",
		Email:       "test@example.com",
		Name:        "tester",
		IsSuperuser: true,
	})

	testCases := []struct {
		name                   string
		request                *http.Request
		sessionService         func(t mockConstructorTestingTNewMockSessionService) *MockSessionService
		memberService          func(t mockConstructorTestingTNewMockMemberService) *MockMemberService
		expectedCode           int
		expectedBody           string
		expectedLoggerMessages []testlogger.Message
	}{
		{
			name: "success",
			request: httptest.NewRequest(
				http.MethodPut,
				"/email",
				strings.NewReader(`{"NewEmail1": "test@example.com", "NewEmail2": "test@example.com"}`),
			).WithContext(contextWithSession),
			sessionService: func(t mockConstructorTestingTNewMockSessionService) *MockSessionService {
				m := NewMockSessionService(t)
				m.
					On("UpdateEmail", mock.Anything, mock.Anything, "test@example.com").
					Return(nil).
					Once()
				return m
			},
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("UpdateEmail", mock.Anything, model.MemberID("42"), "test@example.com").
					Return(nil).
					Once()
				return m
			},
			expectedCode: http.StatusOK,
			expectedBody: `{
				"Status": "Email updated"
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "spaces is trimmed",
			request: httptest.NewRequest(
				http.MethodPut,
				"/email",
				strings.NewReader(`{"NewEmail1": "  test@example.com  ", "NewEmail2": "  test@example.com  "}`),
			).WithContext(contextWithSession),
			sessionService: func(t mockConstructorTestingTNewMockSessionService) *MockSessionService {
				m := NewMockSessionService(t)
				m.
					On("UpdateEmail", mock.Anything, mock.Anything, "test@example.com").
					Return(nil).
					Once()
				return m
			},
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("UpdateEmail", mock.Anything, model.MemberID("42"), "test@example.com").
					Return(nil).
					Once()
				return m
			},
			expectedCode: http.StatusOK,
			expectedBody: `{
				"Status": "Email updated"
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "not authenticated",
			request: httptest.NewRequest(
				http.MethodPut,
				"/email",
				strings.NewReader(`{"NewEmail1": "test@example.com", "NewEmail2": "test@example.com"}`),
			),
			sessionService: NewMockSessionService,
			memberService:  NewMockMemberService,
			expectedCode:   http.StatusUnauthorized,
			expectedBody: `{
				"Status": "Unauthorized",
				"Messages": ["not authenticated"]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "invalid json",
			request: httptest.NewRequest(
				http.MethodPut,
				"/email",
				strings.NewReader(`{`),
			).WithContext(contextWithSession),
			sessionService: NewMockSessionService,
			memberService:  NewMockMemberService,
			expectedCode:   http.StatusBadRequest,
			expectedBody: `{
				"Status": "Bad Request",
				"Messages": ["Invalid JSON"]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "new email1 mismatch",
			request: httptest.NewRequest(
				http.MethodPut,
				"/email",
				strings.NewReader(`{"NewEmail2": "test@example.com"}`),
			).WithContext(contextWithSession),
			sessionService: NewMockSessionService,
			memberService:  NewMockMemberService,
			expectedCode:   http.StatusBadRequest,
			expectedBody: `{
				"Status": "Bad Request",
				"Messages": ["Email must have more than 0 characters.", "Emails don't match."]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "new email2 mismatch",
			request: httptest.NewRequest(
				http.MethodPut,
				"/email",
				strings.NewReader(`{"NewEmail1": "test@example.com"}`),
			).WithContext(contextWithSession),
			sessionService: NewMockSessionService,
			memberService:  NewMockMemberService,
			expectedCode:   http.StatusBadRequest,
			expectedBody: `{
				"Status": "Bad Request",
				"Messages": ["Email must have more than 0 characters.", "Emails don't match."]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "emails don't match",
			request: httptest.NewRequest(
				http.MethodPut,
				"/email",
				strings.NewReader(`{"NewEmail1": "test@example.com", "NewEmail2": "test2@example.com"}`),
			).WithContext(contextWithSession),
			sessionService: NewMockSessionService,
			memberService:  NewMockMemberService,
			expectedCode:   http.StatusBadRequest,
			expectedBody: `{
				"Status": "Bad Request",
				"Messages": ["Emails don't match."]
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "email already exists",
			request: httptest.NewRequest(
				http.MethodPut,
				"/email",
				strings.NewReader(`{"NewEmail1": "test@example.com", "NewEmail2": "test@example.com"}`),
			).WithContext(contextWithSession),
			sessionService: NewMockSessionService,
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("UpdateEmail", mock.Anything, model.MemberID("42"), "test@example.com").
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
			name: "fail update email",
			request: httptest.NewRequest(
				http.MethodPut,
				"/email",
				strings.NewReader(`{"NewEmail1": "test@example.com", "NewEmail2": "test@example.com"}`),
			).WithContext(contextWithSession),
			sessionService: NewMockSessionService,
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("UpdateEmail", mock.Anything, model.MemberID("42"), "test@example.com").
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
					Message: "fail update member email",
					KeyValues: map[string]any{
						"error": "test",
					},
				},
			},
		},
		{
			name: "fail update session",
			request: httptest.NewRequest(
				http.MethodPut,
				"/email",
				strings.NewReader(`{"NewEmail1": "test@example.com", "NewEmail2": "test@example.com"}`),
			).WithContext(contextWithSession),
			sessionService: func(t mockConstructorTestingTNewMockSessionService) *MockSessionService {
				m := NewMockSessionService(t)
				m.
					On("UpdateEmail", mock.Anything, mock.Anything, "test@example.com").
					Return(errors.New("test")).
					Once()
				return m
			},
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("UpdateEmail", mock.Anything, model.MemberID("42"), "test@example.com").
					Return(nil).
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
					Message: "fail update session email",
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
			Setup(router, test.sessionService(t), test.memberService(t))
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
