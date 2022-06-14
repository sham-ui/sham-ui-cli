package sessions

import (
	"database/sql"
	"fmt"
	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"net/http"
	"time"
)

type Store struct {
	store *pgstore.PGStore
}

func (s *Store) GetRawSession(r *http.Request) (*sessions.Session, error) {
	session, err := s.store.Get(r, "cms-session")
	if nil != err {
		return nil, fmt.Errorf("get session from store: %s", err)
	}
	return session, nil
}

func (s *Store) GetSession(r *http.Request) (*Session, error) {
	session, err := s.GetRawSession(r)
	if nil != err {
		return nil, fmt.Errorf("get raw session: %s", err)
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return nil, nil
	}
	id, _ := session.Values["id"].(string)
	name, _ := session.Values["name"].(string)
	email, _ := session.Values["email"].(string)
	isSuperuser, _ := session.Values["is_superuser"].(bool)
	return &Session{
		ID:          id,
		Name:        name,
		Email:       email,
		IsSuperuser: isSuperuser,
		rawSession:  session,
	}, nil
}

func (s *Store) Cleanup(interval time.Duration) (chan<- struct{}, <-chan struct{}) {
	return s.store.Cleanup(interval)
}

func (s *Store) StopCleanup(quit chan<- struct{}, done <-chan struct{}) {
	s.store.StopCleanup(quit, done)
}

func NewStore(db *sql.DB, secret string) (*Store, error) {
	store, err := pgstore.NewPGStoreFromPool(db, []byte(secret))
	if nil != err {
		return nil, fmt.Errorf("can't create pgstore: %s", err)
	}
	return &Store{
		store: store,
	}, nil
}
