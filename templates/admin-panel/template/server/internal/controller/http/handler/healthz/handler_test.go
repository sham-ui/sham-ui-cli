package healthz

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"{{ shortName }}/pkg/asserts"

	"github.com/gorilla/mux"
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
			request:      httptest.NewRequest(http.MethodGet, "http://localhost/healthz", nil),
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := mux.NewRouter()
			Setup(router)
			resp := httptest.NewRecorder()

			// Act
			router.ServeHTTP(resp, test.request)

			// Assert
			asserts.Equals(t, test.expectedBody, resp.Body.String(), "body")
			asserts.Equals(t, test.expectedCode, resp.Code, "code")
		})
	}
}
