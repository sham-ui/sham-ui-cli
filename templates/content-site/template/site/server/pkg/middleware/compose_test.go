package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"site/pkg/asserts"
	"testing"

	"github.com/gorilla/mux"
)

type ctxKey int

func TestCompose(t *testing.T) {
	testCases := []struct {
		name             string
		middlewares      []func(http.Handler) http.Handler
		handler          func(http.ResponseWriter, *http.Request)
		expectedResponse string
	}{
		{
			name: "ordered",
			middlewares: []func(http.Handler) http.Handler{
				func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte("a")) //nolint:errcheck
						next.ServeHTTP(w, r)
					})
				},
				func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte("b")) //nolint:errcheck
						next.ServeHTTP(w, r)
					})
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("c")) //nolint:errcheck
			},
			expectedResponse: "abc",
		},
		{
			name: "request context chained",
			middlewares: []func(http.Handler) http.Handler{
				func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						ctx := context.WithValue(r.Context(), ctxKey(0), "a")
						next.ServeHTTP(w, r.WithContext(ctx))
					})
				},
				func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						ctx := context.WithValue(r.Context(), ctxKey(1), "b")
						next.ServeHTTP(w, r.WithContext(ctx))
					})
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				a := r.Context().Value(ctxKey(0)).(string)
				b := r.Context().Value(ctxKey(1)).(string)
				w.Write([]byte(a + b)) //nolint:errcheck
			},
			expectedResponse: "ab",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := mux.NewRouter()
			router.Use(Compose(test.middlewares...))
			router.PathPrefix("/").HandlerFunc(test.handler)

			rec := httptest.NewRecorder()

			// Act
			router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

			// Assert
			asserts.Equals(t, test.expectedResponse, rec.Body.String())
		})
	}
}
