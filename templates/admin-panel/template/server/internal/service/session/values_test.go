package session

import (
	"{{ shortName }}/pkg/asserts"
	"testing"

	"github.com/gorilla/sessions"
)

func Test_getValue(t *testing.T) {
	testCases := []struct {
		name          string
		session       *sessions.Session
		key           string
		expected      any
		expectedError error
	}{
		{
			name:          "value not found",
			session:       &sessions.Session{}, //nolint:exhaustruct
			key:           "key",
			expected:      "",
			expectedError: newValueNotFoundError("key"),
		},
		{
			name: "value type error",
			session: &sessions.Session{ //nolint:exhaustruct
				Values: map[any]any{
					"key": 42,
				},
			},
			key:           "key",
			expected:      "",
			expectedError: newValueTypeError("key"),
		},
		{
			name: "ok",
			session: &sessions.Session{ //nolint:exhaustruct
				Values: map[any]any{
					"key": "value",
				},
			},
			key:      "key",
			expected: "value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := getValue[string](tc.session, tc.key)
			asserts.Equals(t, tc.expected, value)
			asserts.ErrorsEqual(t, tc.expectedError, err)
		})
	}
}
