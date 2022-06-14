package sessions

import "github.com/gorilla/sessions"

type Session struct {
	ID          string
	Email       string
	Name        string
	IsSuperuser bool

	rawSession *sessions.Session
}

func (s *Session) GetRawSession() *sessions.Session {
	return s.rawSession
}
