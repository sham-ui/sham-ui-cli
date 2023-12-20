package ssr

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"site/pkg/asserts"
	"site/pkg/logger"

	"github.com/go-logr/logr/testr"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ServeHTTP(t *testing.T) {
	testCases := []struct {
		name         string
		render       func(t mockConstructorTestingTNewMockServerSideRender) *MockServerSideRender
		request      *http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name: "success",
			render: func(t mockConstructorTestingTNewMockServerSideRender) *MockServerSideRender {
				u, _ := url.ParseRequestURI("http://localhost/page-name")
				c := []*http.Cookie\{{
					Name:  "foo",
					Value: "bar",
				}}
				m := NewMockServerSideRender(t)
				m.On("Render", mock.Anything, u, c).Return([]byte("html"), nil)
				return m
			},
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "http://localhost/page-name", nil)
				r.AddCookie(&http.Cookie{
					Name:  "foo",
					Value: "bar",
				})
				return r
			}(),
			expectedCode: http.StatusOK,
			expectedBody: "html",
		},
		{
			name: "fail",
			render: func(t mockConstructorTestingTNewMockServerSideRender) *MockServerSideRender {
				u, _ := url.ParseRequestURI("http://localhost/page-name")
				m := NewMockServerSideRender(t)
				m.On("Render", mock.Anything, u, []*http.Cookie{}).Return(nil, errors.New("test error"))
				return m
			},
			request:      httptest.NewRequest(http.MethodGet, "http://localhost/page-name", nil),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "{\"Status\":\"Internal Server Error\",\"Messages\":[\"internal server error\"]}\n",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := mux.NewRouter()
			router.PathPrefix("/").Handler(NewHandler(test.render(t)))
			ctx := logger.Save(test.request.Context(), testr.New(t))
			request := test.request.WithContext(ctx)
			resp := httptest.NewRecorder()

			// Act
			router.ServeHTTP(resp, request)

			// Assert
			asserts.Equals(t, test.expectedBody, resp.Body.String(), "body")
			asserts.Equals(t, test.expectedCode, resp.Code, "code")
		})
	}
}
