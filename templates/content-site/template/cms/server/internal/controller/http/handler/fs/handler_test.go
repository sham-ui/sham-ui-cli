package fs

import (
	"cms/internal/controller/http/request"
	"cms/internal/model"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cms/pkg/asserts"
	"cms/pkg/logger"

	"github.com/go-logr/logr/testr"
	"github.com/gorilla/mux"
)

func TestHandler_ServeHTTP(t *testing.T) {
	fsys := os.DirFS("testdata")
	htmlHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("html"))
	})

	testCases := []struct {
		name         string
		request      *http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name:         "file.txt",
			request:      httptest.NewRequest(http.MethodGet, "/file.txt", nil),
			expectedCode: http.StatusOK,
			expectedBody: "text file content",
		},
		{
			name:         "texts/file.txt",
			request:      httptest.NewRequest(http.MethodGet, "/texts/file.txt", nil),
			expectedCode: http.StatusOK,
			expectedBody: "text file content from dir",
		},
		{
			name:         "texts dir",
			request:      httptest.NewRequest(http.MethodGet, "/texts", nil),
			expectedCode: http.StatusOK,
			expectedBody: "html",
		},
		{
			name:         "not found",
			request:      httptest.NewRequest(http.MethodGet, "/not-found", nil),
			expectedCode: http.StatusOK,
			expectedBody: "html",
		},
		{
			name:         "index.html",
			request:      httptest.NewRequest(http.MethodGet, "/index.html", nil),
			expectedCode: http.StatusOK,
			expectedBody: "html",
		},
		{
			name:         "root path",
			request:      httptest.NewRequest(http.MethodGet, "/", nil),
			expectedCode: http.StatusOK,
			expectedBody: "html",
		},
		{
			name: "superuser file",
			request: func() *http.Request {
				ctx := request.SaveSessionToContext(
					context.Background(),
					&model.Session{ //nolint:exhaustruct
						IsSuperuser: true,
					},
				)
				return httptest.NewRequest(http.MethodGet, "/su_file.txt", nil).
					WithContext(ctx)
			}(),
			expectedCode: http.StatusOK,
			expectedBody: "Superuser file",
		},
		{
			name:         "superuser file unauthorized",
			request:      httptest.NewRequest(http.MethodGet, "/su_file.txt", nil),
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Unauthorized\n",
		},
		{
			name: "superuser file forbidden",
			request: func() *http.Request {
				ctx := request.SaveSessionToContext(
					context.Background(),
					&model.Session{ //nolint:exhaustruct
						IsSuperuser: false,
					},
				)
				return httptest.NewRequest(http.MethodGet, "/su_file.txt", nil).
					WithContext(ctx)
			}(),
			expectedCode: http.StatusForbidden,
			expectedBody: "Forbidden\n",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := mux.NewRouter()
			Setup(router, fsys, htmlHandler)
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
