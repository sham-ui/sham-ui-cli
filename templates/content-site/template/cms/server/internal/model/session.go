package model

import "errors"

type Session struct {
	MemberID    MemberID
	Email       string
	Name        string
	IsSuperuser bool
}

var (
	ErrSessionNotAuthenticated = errors.New("session not authenticated")
	ErrSessionNotSuperuser     = errors.New("session not superuser")
)
