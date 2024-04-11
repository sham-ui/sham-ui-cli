package model

import "errors"

type (
	CategoryID   string
	CategorySlug string

	Category struct {
		ID   CategoryID
		Slug CategorySlug
		Name string
	}
)

var (
	ErrCategoryNotFound          = errors.New("category not found")
	ErrCategoryNameAlreadyExists = errors.New("category name already exists")
	ErrCategorySlugAlreadyExists = errors.New("category slug already exists")
)
