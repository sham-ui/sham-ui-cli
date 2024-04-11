package server_info

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"{{ shortName }}/pkg/asserts"
	"{{ shortName }}/pkg/logger"
	"{{ shortName }}/pkg/logger/testlogger"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
)

//go:embed testdata/*
var testdata embed.FS

func TestHandler_ServeHTTP(t *testing.T) {
	// Arrange
	startTime := time.Date(2024, time.January, 17, 20, 31, 14, 0, time.UTC)
	testFS, err := fs.Sub(testdata, "testdata")
	asserts.NoError(t, err)
	router := mux.NewRouter()
	Setup(router, testFS, startTime)
	log := testlogger.NewLogger()
	req := httptest.NewRequest(http.MethodGet, "/server-info", nil).WithContext(
		logger.Save(context.Background(), logr.New(log)),
	)
	resp := httptest.NewRecorder()

	// Action
	router.ServeHTTP(resp, req)

	// Assert
	asserts.Equals(t, http.StatusOK, resp.Code, "code")
	asserts.Equals(t, []testlogger.Message{}, log.Messages, "logger")
	decodedResponse := make(map[string]any)
	err = json.NewDecoder(resp.Body).Decode(&decodedResponse)
	asserts.NoError(t, err)
	keys := make([]string, 0, len(decodedResponse))
	for k := range decodedResponse {
		keys = append(keys, k)
	}
	asserts.EqualsIgnoreOrder(t, []string{"Host", "Runtime", "Files"}, keys, "response keys")
	asserts.Equals(t, []any{
		map[string]any{
			"Name": "1.txt",
			"Size": float64(12),
		},
		map[string]any{
			"Name": "2.png",
			"Size": float64(0),
		},
		map[string]any{
			"Name": "sub-dir/1.txt",
			"Size": float64(23),
		},
	}, decodedResponse["Files"], "files")
}
