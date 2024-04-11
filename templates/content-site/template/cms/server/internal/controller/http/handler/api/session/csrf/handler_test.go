package csrf

import (
	"cms/pkg/asserts"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHandler_ServeHTTP(t *testing.T) {
	headerKey := "X-CSRF-Token"
	contextKey := "gorilla.csrf.Token"

	testCases := []struct {
		name                   string
		request                *http.Request
		expectedResponseHeader string
	}{
		{
			name: "success",
			request: httptest.NewRequest(http.MethodGet, "/csrftoken", nil).WithContext(
				context.WithValue(context.Background(), contextKey, "token"), //nolint:staticcheck
			),
			expectedResponseHeader: "token",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := mux.NewRouter()
			Setup(router, headerKey)
			resp := httptest.NewRecorder()

			// Action
			router.ServeHTTP(resp, test.request)

			// Assert
			asserts.Equals(t, test.expectedResponseHeader, resp.Header().Get(headerKey))
		})
	}
}
