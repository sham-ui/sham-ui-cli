package recovery

import (
	"net/http"
	"net/http/httptest"
	"{{ shortName }}/pkg/asserts"
	"{{ shortName }}/pkg/logger/testlogger"
	"testing"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
)

func TestRecovery_ServeHTTP(t *testing.T) {
	// Arrange
	log := testlogger.NewLogger()

	router := mux.NewRouter()
	router.Use(New(logr.New(log)).Middleware)
	router.Handle("/{path}", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		panic("test")
	}))

	req := httptest.NewRequest(http.MethodGet, "/foo?bar=42", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	asserts.Equals(t, http.StatusInternalServerError, resp.Code, "code")
	asserts.JSONEquals(t,
		`{"Status":"Internal Server Error","Messages":["internal server error"]}`,
		resp.Body.String(),
		"body",
	)
	asserts.Equals(t, 1, len(log.Messages), "log count")
	logMsg := log.Messages[0]
	asserts.Equals(t, 0, logMsg.Level, "log level")
	asserts.Equals(t, "panic recovered", logMsg.Message, "log message")
	asserts.Equals(t, 2, len(logMsg.KeyValues), "log key values count")
	asserts.Equals(t, "http handler panic: GET /foo?bar=42: test", logMsg.KeyValues["error"], "log error")
	asserts.Equals(t, true, logMsg.HasKey("stack"), "log stack")
}
