package model

import "errors"

type (
	TagID   string
	TagSlug string

	Tag struct {
		ID   TagID
		Slug TagSlug
		Name string
	}
)

var (
	ErrTagNotFound          = errors.New("tag not found")
	ErrTagNameAlreadyExists = errors.New("tag name already exists")
	ErrTagSlugAlreadyExists = errors.New("tag slug already exists")
)
