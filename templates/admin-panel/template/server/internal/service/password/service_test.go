package password

import (
	"context"
	"errors"
	"{{ shortName }}/internal/model"
	"{{ shortName }}/pkg/asserts"
	"testing"
)

func TestService_Hash(t *testing.T) {
	// Arrange
	ctx := context.Background()
	srv := New()
	const password = "test"

	// Act
	hashed, err := srv.Hash(ctx, password)
	asserts.NoError(t, err)

	// Assert
	err = srv.Compare(ctx, hashed, password)
	asserts.NoError(t, err)
}

func TestService_Compare(t *testing.T) {
	testCases := []struct {
		name          string
		hashed        model.MemberHashedPassword
		raw           string
		expectedError error
	}{
		{
			name:   "success",
			hashed: model.MemberHashedPassword("$2a$14$3StvQXPkuXnx75TnxkQvqu25yJnuyg63FmT8E.9exeMZONJa6ljCK"),
			raw:    "test",
		},
		{
			name:   "invalid hash",
			hashed: model.MemberHashedPassword("123"),
			raw:    "test",
			expectedError: errors.New(
				"compare hash and password: crypto/bcrypt: hashedSecret too short to be a bcrypted password",
			),
		},
		{
			name:   "not matching",
			hashed: model.MemberHashedPassword("$2a$14$3StvQXPkuXnx75TnxkQvqu25yJnuyg63FmT8E.9exeMZONJa6ljCK"),
			raw:    "wrong",
			expectedError: errors.New(
				"compare hash and password: crypto/bcrypt: hashedPassword is not the hash of the given password",
			),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			srv := New()

			// Act
			err := srv.Compare(context.Background(), test.hashed, test.raw)

			// Assert
			asserts.ErrorsEqual(t, test.expectedError, err)
		})
	}
}
