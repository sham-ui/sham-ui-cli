package validation

import (
	"testing"

	"cms/pkg/asserts"
)

func TestValidation(t *testing.T) {
	// Arrange
	validation := New()

	// Act
	validation.AddErrors("error1")
	validation.AddErrors("error2")

	// Assert
	asserts.Equals(t, true, validation.HasErrors())
	asserts.Equals(t, []string{"error1", "error2"}, validation.Errors)
}

func TestValidation_AddErrors(t *testing.T) {
	testCases := []struct {
		name     string
		errors   []string
		hasError bool
	}{
		{
			name:     "empty",
			errors:   []string{},
			hasError: false,
		},
		{
			name:     "single",
			errors:   []string{"error1"},
			hasError: true,
		},
		{
			name:     "multiple",
			errors:   []string{"error1", "error2"},
			hasError: true,
		},
		{
			name:     "nil",
			errors:   nil,
			hasError: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			validation := New()
			validation.AddErrors(test.errors...)
			asserts.Equals(t, test.hasError, validation.HasErrors())
		})
	}
}
