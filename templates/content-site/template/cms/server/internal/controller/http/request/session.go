package request

import (
	"cms/internal/model"
	"context"
)

type contextKeySession struct{}

func SessionFromContext(ctx context.Context) (*model.Session, bool) {
	sess, ok := ctx.Value(contextKeySession{}).(*model.Session)
	return sess, ok
}

func SaveSessionToContext(ctx context.Context, session *model.Session) context.Context {
	return context.WithValue(ctx, contextKeySession{}, session)
}
