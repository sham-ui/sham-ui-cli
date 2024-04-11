package metrics

import (
	"cms/pkg/asserts"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestSetup(t *testing.T) {
	// Arrange
	router := mux.NewRouter()
	Setup(router)

	// Act
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "http://localhost/metrics", nil))

	// Assert
	asserts.Equals(t, http.StatusOK, resp.Code, "code")
}
