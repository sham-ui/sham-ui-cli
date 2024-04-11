package request

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"{{ shortName }}/pkg/asserts"
	"strings"
	"testing"
)

func TestDecodeJSONt(t *testing.T) {
	type data struct {
		Key string `json:"key"`
	}
	testCases := []struct {
		name          string
		request       *http.Request
		expected      *data
		expectedError error
	}{
		{
			name:    "valid request",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"key": "value"}`)),
			expected: &data{
				Key: "value",
			},
			expectedError: nil,
		},
		{
			name:          "invalid request",
			request:       httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"key":}`)),
			expected:      nil,
			expectedError: errors.New("decode JSON: invalid character '}' looking for beginning of value"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Act
			res, err := DecodeJSON[data](test.request)

			// Assert
			asserts.Equals(t, test.expected, res, "data")
			asserts.ErrorsEqual(t, test.expectedError, err)
		})
	}
}
