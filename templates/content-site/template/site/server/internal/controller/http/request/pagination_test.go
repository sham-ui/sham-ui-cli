package request

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"site/pkg/asserts"
)

func TestPaginationParams_Validation(t *testing.T) {
	testCases := []struct {
		name                   string
		request                *http.Request
		expected               *PaginationParams
		expectedRequestValid   bool
		expectedResponseStatus int
		expectedResponseBody   string
	}{
		{
			name:                   "success",
			request:                httptest.NewRequest(http.MethodGet, "/?offset=5&limit=10", nil),
			expected:               &PaginationParams{Offset: 5, Limit: 10},
			expectedRequestValid:   true,
			expectedResponseStatus: http.StatusOK,
		},
		{
			name:                   "default offset",
			request:                httptest.NewRequest(http.MethodGet, "/?limit=10", nil),
			expected:               &PaginationParams{Offset: DefaultOffset, Limit: 10},
			expectedRequestValid:   true,
			expectedResponseStatus: http.StatusOK,
		},
		{
			name:                   "default limit",
			request:                httptest.NewRequest(http.MethodGet, "/?offset=5", nil),
			expected:               &PaginationParams{Offset: 5, Limit: DefaultLimit},
			expectedRequestValid:   true,
			expectedResponseStatus: http.StatusOK,
		},
		{
			name:                   "invalid offset",
			request:                httptest.NewRequest(http.MethodGet, "/?offset=abc&limit=10", nil),
			expected:               &PaginationParams{Offset: DefaultOffset, Limit: 10},
			expectedRequestValid:   true,
			expectedResponseStatus: http.StatusOK,
		},
		{
			name:                   "invalid limit",
			request:                httptest.NewRequest(http.MethodGet, "/?offset=5&limit=abc", nil),
			expected:               &PaginationParams{Offset: 5, Limit: DefaultLimit},
			expectedRequestValid:   true,
			expectedResponseStatus: http.StatusOK,
		},
		{
			name:                   "limit must be greater than 0",
			request:                httptest.NewRequest(http.MethodGet, "/?offset=5&limit=-1", nil),
			expected:               nil,
			expectedRequestValid:   false,
			expectedResponseStatus: http.StatusBadRequest,
			expectedResponseBody:   `{"Messages": ["limit must be greater than 0"], "Status": "Bad Request"}`,
		},
		{
			name:                   "offset must be greater or equal than 0",
			request:                httptest.NewRequest(http.MethodGet, "/?offset=-1&limit=10", nil),
			expected:               nil,
			expectedRequestValid:   false,
			expectedResponseStatus: http.StatusBadRequest,
			expectedResponseBody:   `{"Messages": ["offset must be greater or equal than 0"], "Status": "Bad Request"}`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			resp := httptest.NewRecorder()

			// Act
			p, valid := ExtractPaginationParams(resp, test.request)

			// Assert
			asserts.Equals(t, test.expected, p)
			asserts.Equals(t, test.expectedRequestValid, valid)
			asserts.Equals(t, test.expectedResponseStatus, resp.Code)
			if test.expectedResponseBody == "" {
				asserts.Equals(t, "", resp.Body.String())
			} else {
				asserts.JSONEquals(t, test.expectedResponseBody, resp.Body.String())
			}
		})
	}
}
