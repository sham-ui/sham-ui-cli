package category

import "cms/internal/model"

type category struct {
	ID   string
	Name string
	Slug string
}

func (c category) toModel() model.Category {
	return model.Category{
		ID:   model.CategoryID(c.ID),
		Name: c.Name,
		Slug: model.CategorySlug(c.Slug),
	}
}
