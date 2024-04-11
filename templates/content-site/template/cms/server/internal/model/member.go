package model

import "errors"

type (
	MemberID             string
	MemberHashedPassword []byte

	Member struct {
		ID          MemberID
		Email       string
		Name        string
		IsSuperuser bool
	}

	MemberWithPassword struct {
		Member
		HashedPassword MemberHashedPassword
	}
)

var (
	ErrMemberNotFound           = errors.New("member not found")
	ErrMemberEmailAlreadyExists = errors.New("member email already exists")
)
