package cms

import (
	"site/internal/external_api/cms/proto"
	"site/internal/model"
)

func articleListItemsToModel(articles []*proto.ArticleListItem) []model.ShortArticle {
	result := make([]model.ShortArticle, len(articles))
	for i, article := range articles {
		result[i] = model.ShortArticle{
			Title: article.GetTitle(),
			Slug:  article.GetSlug(),
			Category: model.Category{
				Name: article.GetCategory().Name,
				Slug: article.GetCategory().Slug,
			},
			ShortContent: article.GetContent(),
			PublishedAt:  article.GetPublishedAt().AsTime(),
		}
	}
	return result
}
