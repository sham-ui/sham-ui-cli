package request

import (
	"cms/internal/model"
	"cms/pkg/asserts"
	"context"
	"testing"
)

func Test_Session(t *testing.T) {
	// Arrange
	ctx := context.Background()
	sess := &model.Session{
		MemberID:    "42",
		Email:       "test@example.com",
		Name:        "tester",
		IsSuperuser: true,
	}

	// Action
	ctx = SaveSessionToContext(ctx, sess)

	// Assert
	sessFromCtx, ok := SessionFromContext(ctx)
	asserts.Equals(t, sess, sessFromCtx, "session")
	asserts.Equals(t, true, ok, "ok")
}
