package model

import "time"

type Article struct {
	ShortArticle
	Tags    []Tag
	Content string
}

type ShortArticle struct {
	Title        string
	Slug         string
	Category     Category
	ShortContent string
	PublishedAt  time.Time
}

type ArticleNotFoundError struct {
	Slug string
}

func (e ArticleNotFoundError) Error() string {
	return "article not found: " + e.Slug
}

func (e ArticleNotFoundError) Is(target error) bool {
	//nolint:errorlint
	_, ok := target.(ArticleNotFoundError)
	return ok
}

func NewArticleNotFoundError(slug string) error {
	return ArticleNotFoundError{
		Slug: slug,
	}
}
