package model

import (
	"errors"
	"time"
)

type (
	ArticleID   string
	ArticleSlug string

	Article struct {
		ID          ArticleID
		Slug        ArticleSlug
		Title       string
		CategoryID  CategoryID
		ShortBody   string
		Body        string
		PublishedAt time.Time
	}

	ArticleShortInfo struct {
		ID          ArticleID
		Slug        ArticleSlug
		Title       string
		CategoryID  CategoryID
		PublishedAt time.Time
	}

	ArticleShortInfoWithCategory struct {
		ID          ArticleID
		Slug        ArticleSlug
		Title       string
		Category    Category
		ShortBody   string
		PublishedAt time.Time
	}
)

var (
	ErrArticleNotFound           = errors.New("article not found")
	ErrArticleTitleAlreadyExists = errors.New("article title already exists")
	ErrArticleSlugAlreadyExists  = errors.New("article slug already exists")
)
