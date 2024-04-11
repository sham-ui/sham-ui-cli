package session

import (
	"cms/internal/controller/http/request"
	"cms/internal/model"
	"cms/pkg/logger"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Service interface {
	Get(r *http.Request) (*model.Session, error)
}

type session struct {
	srv Service
}

func (sess *session) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			memberSess, err := sess.srv.Get(r)
			if err == nil {
				ctx := r.Context()

				// Save session to trace context
				span := trace.SpanFromContext(ctx)
				if span.SpanContext().IsValid() {
					span.SetAttributes(
						attribute.String("member_id", string(memberSess.MemberID)),
					)
				}

				// Save session to logger context
				log := logger.Load(ctx).WithValues(
					"member_id", string(memberSess.MemberID),
				)
				ctx = logger.Save(r.Context(), log)

				// Save session to request context
				ctx = request.SaveSessionToContext(ctx, memberSess)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(rw, r)
		},
	)
}

func New(srv Service) *session {
	return &session{
		srv: srv,
	}
}
