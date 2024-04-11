package session

import (
	"cms/internal/model"
	"cms/pkg/postgres"
	"cms/pkg/tracing"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/antonlindstrom/pgstore"
)

const (
	scopeName = "service.session"

	sessionName = "cms-session"

	sessionValueAuthenticated = "authenticated"
	sessionValueMemberID      = "id"
	sessionValueEmail         = "email"
	sessionValueName          = "name"
	sessionValueIsSuperuser   = "is_superuser"
)

type service struct {
	store *pgstore.PGStore

	quit chan<- struct{}
	done <-chan struct{}
}

func (s *service) String() string {
	return scopeName
}

func (s *service) GracefulShutdown(_ context.Context) error {
	s.store.StopCleanup(s.quit, s.done)
	return nil
}

func (s *service) Get(r *http.Request) (*model.Session, error) {
	const op = "Get"

	_, span := tracing.StartSpan(r.Context(), scopeName, op)
	defer span.End()

	raw, err := s.store.Get(r, sessionName)
	if err != nil {
		return nil, fmt.Errorf("get session from store: %w", err)
	}
	if auth, err := getValue[bool](raw, sessionValueAuthenticated); err != nil || !auth {
		return nil, model.ErrSessionNotAuthenticated
	}
	session := &model.Session{} //nolint:exhaustruct
	memberID, err := getValue[string](raw, sessionValueMemberID)
	if err != nil {
		return nil, fmt.Errorf("get member id: %w", err)
	}
	session.MemberID = model.MemberID(memberID)
	if session.Name, err = getValue[string](raw, sessionValueName); err != nil {
		return nil, fmt.Errorf("get name: %w", err)
	}
	if session.Email, err = getValue[string](raw, sessionValueEmail); err != nil {
		return nil, fmt.Errorf("get email: %w", err)
	}
	if session.IsSuperuser, err = getValue[bool](raw, sessionValueIsSuperuser); err != nil {
		return nil, fmt.Errorf("get is_superuser: %w", err)
	}
	return session, nil
}

func (s *service) Create(rw http.ResponseWriter, r *http.Request, member *model.Member) error {
	const op = "Create"

	_, span := tracing.StartSpan(r.Context(), scopeName, op)
	defer span.End()

	raw, err := s.store.New(r, sessionName)
	if err != nil {
		return fmt.Errorf("can't create new session: %w", err)
	}
	raw.Values[sessionValueAuthenticated] = true
	raw.Values[sessionValueMemberID] = string(member.ID)
	raw.Values[sessionValueName] = member.Name
	raw.Values[sessionValueEmail] = member.Email
	raw.Values[sessionValueIsSuperuser] = member.IsSuperuser
	if err := raw.Save(r, rw); err != nil {
		return fmt.Errorf("can't save session: %w", err)
	}
	return nil
}

func (s *service) Delete(rw http.ResponseWriter, r *http.Request) error {
	const op = "Delete"

	_, span := tracing.StartSpan(r.Context(), scopeName, op)
	defer span.End()

	raw, err := s.store.Get(r, sessionName)
	if err != nil {
		return fmt.Errorf("get session from store: %w", err)
	}
	raw.Values[sessionValueAuthenticated] = false
	raw.Options.MaxAge = -1
	if err := raw.Save(r, rw); err != nil {
		return fmt.Errorf("can't save session: %w", err)
	}
	return nil
}

func (s *service) UpdateEmail(rw http.ResponseWriter, r *http.Request, email string) error {
	const op = "UpdateEmail"
	_, span := tracing.StartSpan(r.Context(), scopeName, op)
	defer span.End()
	if err := s.updateValue(rw, r, sessionValueEmail, email); err != nil {
		return fmt.Errorf("update email: %w", err)
	}
	return nil
}

func (s *service) UpdateName(rw http.ResponseWriter, r *http.Request, name string) error {
	const op = "UpdateName"
	_, span := tracing.StartSpan(r.Context(), scopeName, op)
	defer span.End()
	if err := s.updateValue(rw, r, sessionValueName, name); err != nil {
		return fmt.Errorf("update name: %w", err)
	}
	return nil
}

func (s *service) updateValue(rw http.ResponseWriter, r *http.Request, field string, value any) error {
	raw, err := s.store.Get(r, sessionName)
	if err != nil {
		return fmt.Errorf("get session from store: %w", err)
	}
	raw.Values[field] = value
	if err := raw.Save(r, rw); err != nil {
		return fmt.Errorf("can't save session: %w", err)
	}
	return nil
}

func New(
	db *postgres.Database,
	secret string,
	domain string,
	maxAge time.Duration,
) (*service, error) {
	store, err := pgstore.NewPGStoreFromPool(db.Postgres(), []byte(secret))
	if err != nil {
		return nil, fmt.Errorf("can't create pgstore: %w", err)
	}
	store.MaxAge(int(maxAge.Seconds()))
	store.Options.Domain = domain
	quit, done := store.Cleanup(maxAge)
	return &service{
		store: store,
		quit:  quit,
		done:  done,
	}, nil
}
