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
		memberService          func(t mockConstructorTestingTNewMockMemberService) *MockMemberService
		expectedCode           int
		expectedBody           string
		expectedLoggerMessages []testlogger.Message
	}{
		{
			name: "success",
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("Find", mock.Anything, int64(10), int64(30)).
					Return([]model.Member{
						{
							ID:          "42",
							Email:       "test1@example.com",
							Name:        "tester-1",
							IsSuperuser: true,
						},
						{
							ID:          "43",
							Email:       "test2@example.com",
							Name:        "tester-2",
							IsSuperuser: false,
						},
					}, nil).
					Once()
				m.
					On("Total", mock.Anything).
					Return(12, nil).
					Once()
				return m
			},
			request:      httptest.NewRequest(http.MethodGet, "/?offset=10&limit=30", nil),
			expectedCode: http.StatusOK,
			expectedBody: `{
			  "members": [
				{
				  "ID": "42",
				  "Name": "tester-1",
				  "Email": "test1@example.com",
				  "IsSuperuser": true
				},
				{
				  "ID": "43",
				  "Name": "tester-2",
				  "Email": "test2@example.com",
				  "IsSuperuser": false
				}
			  ],
			  "meta": {
				"offset": 10,
				"limit": 30,
				"total": 12
			  }
			}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:                   "not valid params",
			memberService:          NewMockMemberService,
			request:                httptest.NewRequest(http.MethodGet, "/?limit=-1", nil),
			expectedCode:           http.StatusBadRequest,
			expectedBody:           `{"Status":"Bad Request","Messages":["limit must be greater than 0"]}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name: "find error",
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("Find", mock.Anything, int64(10), int64(30)).
					Return(nil, errors.New("test")).
					Once()
				return m
			},
			request:      httptest.NewRequest(http.MethodGet, "/?offset=10&limit=30", nil),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"Status":"Internal Server Error","Messages":["internal server error"]}`,
			expectedLoggerMessages: []testlogger.Message{
				{
					Message:   "failed to find members",
					KeyValues: map[string]interface{}{"error": "test"},
				},
			},
		},
		{
			name: "total error",
			memberService: func(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
				m := NewMockMemberService(t)
				m.
					On("Find", mock.Anything, int64(10), int64(30)).
					Return(nil, nil).
					Once()
				m.
					On("Total", mock.Anything).
					Return(0, errors.New("test")).
					Once()
				return m
			},
			request:      httptest.NewRequest(http.MethodGet, "/?offset=10&limit=30", nil),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"Status":"Internal Server Error","Messages":["internal server error"]}`,
			expectedLoggerMessages: []testlogger.Message{
				{
					Message:   "failed to count members",
					KeyValues: map[string]interface{}{"error": "test"},
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
