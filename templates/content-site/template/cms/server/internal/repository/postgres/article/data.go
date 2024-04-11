package article

import (
	"cms/internal/model"
	"time"
)

type article struct {
	ID           string
	Title        string
	Slug         string
	CategoryID   string
	CategoryName string
	CategorySlug string
	ShortBody    string
	Body         string
	PublishedAt  time.Time
}

func (a *article) toModel() model.Article {
	return model.Article{
		ID:          model.ArticleID(a.ID),
		Title:       a.Title,
		Slug:        model.ArticleSlug(a.Slug),
		CategoryID:  model.CategoryID(a.CategoryID),
		ShortBody:   a.ShortBody,
		Body:        a.Body,
		PublishedAt: a.PublishedAt,
	}
}

func (a *article) toModelShortInfo() model.ArticleShortInfo {
	return model.ArticleShortInfo{
		ID:          model.ArticleID(a.ID),
		Slug:        model.ArticleSlug(a.Slug),
		Title:       a.Title,
		CategoryID:  model.CategoryID(a.CategoryID),
		PublishedAt: a.PublishedAt,
	}
}

func (a *article) toModelWithCategory() model.ArticleShortInfoWithCategory {
	return model.ArticleShortInfoWithCategory{
		ID:    model.ArticleID(a.ID),
		Slug:  model.ArticleSlug(a.Slug),
		Title: a.Title,
		Category: model.Category{
			ID:   model.CategoryID(a.CategoryID),
			Slug: model.CategorySlug(a.CategorySlug),
			Name: a.CategoryName,
		},
		ShortBody:   a.ShortBody,
		PublishedAt: a.PublishedAt,
	}
}
