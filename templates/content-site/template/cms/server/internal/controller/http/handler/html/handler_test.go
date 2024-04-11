package html

import (
	"cms/pkg/asserts"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ServeHTTP(t *testing.T) {
	testCases := []struct {
		name         string
		request      *http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name:         "success",
			request:      httptest.NewRequest(http.MethodGet, "/foo", nil),
			expectedCode: http.StatusOK,
			expectedBody: "test",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			h := New([]byte("test"))
			resp := httptest.NewRecorder()

			// Action
			h.ServeHTTP(resp, test.request)

			// Assert
			asserts.Equals(t, test.expectedCode, resp.Code, "code")
			asserts.Equals(t, test.expectedBody, resp.Body.String(), "body")
		})
	}
}
