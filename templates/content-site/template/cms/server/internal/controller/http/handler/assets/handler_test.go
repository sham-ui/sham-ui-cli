package assets

import (
	"cms/pkg/asserts"
	"cms/pkg/logger"
	"cms/pkg/logger/testlogger"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
)

type (
	rawString string
	json      string
)

func TestHandler_ServeHTTP(t *testing.T) {
	testCases := []struct {
		name         string
		request      *http.Request
		fs           fs.FS
		expectedCode int
		expectedBody any
	}{
		{
			name:    "success",
			request: httptest.NewRequest(http.MethodGet, "/test.txt", nil),
			fs: fstest.MapFS{
				"test.txt": {
					Data: []byte("Test"),
				},
			},
			expectedCode: http.StatusOK,
			expectedBody: rawString("Test"),
		},
		{
			name:         "not found",
			request:      httptest.NewRequest(http.MethodGet, "/test.txt", nil),
			fs:           fstest.MapFS{},
			expectedCode: http.StatusNotFound,
			expectedBody: rawString("404 page not found\n"),
		},
		{
			name:         "missing extension",
			request:      httptest.NewRequest(http.MethodGet, "/test", nil),
			fs:           fstest.MapFS{},
			expectedCode: http.StatusBadRequest,
			expectedBody: json(`{"Status":"Bad Request","Messages":["Empty extension"]}`),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := mux.NewRouter()
			Setup(router, test.fs)
			log := testlogger.NewLogger()
			ctx := logger.Save(test.request.Context(), logr.New(log))
			req := test.request.WithContext(ctx)
			resp := httptest.NewRecorder()

			// Action
			router.ServeHTTP(resp, req)

			// Assert
			asserts.Equals(t, test.expectedCode, resp.Code, "code")
			switch expectedBody := test.expectedBody.(type) {
			case rawString:
				asserts.Equals(t, string(expectedBody), resp.Body.String(), "body")
			case json:
				asserts.JSONEquals(t, string(expectedBody), resp.Body.String(), "body")
			default:
				t.Fatalf("unknown type: %t", expectedBody)
			}
		})
	}
}
