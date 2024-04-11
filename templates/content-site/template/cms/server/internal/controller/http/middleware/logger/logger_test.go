package logger

import (
	"cms/internal/controller/http/response"
	"cms/pkg/asserts"
	"cms/pkg/logger/testlogger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
)

func TestLogger_Middleware(t *testing.T) {
	// Arrange
	log := testlogger.NewLogger()

	router := mux.NewRouter()
	router.Use(New(logr.New(log)).Middleware)
	router.Handle("/{path}", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		response.JSON(rw, r, http.StatusCreated, map[string]any{"key": "value"})
	}))

	req := httptest.NewRequest(http.MethodGet, "/foo?bar=42", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	asserts.Equals(t, http.StatusCreated, resp.Code, "code")
	asserts.JSONEquals(t, `{"key":"value"}`, resp.Body.String(), "body")
	asserts.Equals(t, 1, len(log.Messages), "log count")
	logMsg := log.Messages[0]
	asserts.Equals(t, 0, logMsg.Level, "log level")
	asserts.Equals(t, "http request", logMsg.Message, "log message")
	asserts.EqualsIgnoreOrder(t,
		[]string{
			"start",
			"method",
			"hostname",
			"request_uri",
			"user_agent",
			"remote_addr",
			"status",
			"duration",
			"route",
		},
		logMsg.Keys(),
		"extra keys",
	)
	asserts.Equals(t, "GET", logMsg.KeyValues["method"], "method")
	asserts.Equals(t, "example.com", logMsg.KeyValues["hostname"], "hostname")
	asserts.Equals(t, "/foo?bar=42", logMsg.KeyValues["request_uri"], "request_uri")
}
