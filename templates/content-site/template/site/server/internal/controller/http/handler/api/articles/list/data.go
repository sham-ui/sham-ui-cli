package list

import "site/internal/model"

type (
	Article struct {
		Title     string   `json:"title"`
		Slug      string   `json:"slug"`
		Category  Category `json:"category"`
		Content   string   `json:"content"`
		CreatedAt string   `json:"createdAt"`
	}
	Category struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	Meta struct {
		Offset int64 `json:"offset"`
		Limit  int64 `json:"limit"`
		Total  int64 `json:"total"`
	}
)

func ArticlesFromModel(modelArticles []model.ShortArticle) []Article {
	articles := make([]Article, len(modelArticles))
	for i, item := range modelArticles {
		articles[i] = Article{
			Title: item.Title,
			Slug:  item.Slug,
			Category: Category{
				Name: item.Category.Name,
				Slug: item.Category.Slug,
			},
			Content:   item.ShortContent,
			CreatedAt: item.PublishedAt.Format("2006-01-02 15:04:05.999999999 -0700 MST"),
		}
	}
	return articles
}
