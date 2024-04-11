package list

import (
	"cms/internal/controller/grpc/proto"
	"cms/internal/model"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ArticleListItemFromModel(items []model.ArticleShortInfoWithCategory) []*proto.ArticleListItem {
	res := make([]*proto.ArticleListItem, len(items))
	for i, item := range items {
		res[i] = &proto.ArticleListItem{
			Title: item.Title,
			Slug:  string(item.Slug),
			Category: &proto.Category{
				Name: item.Category.Name,
				Slug: string(item.Category.Slug),
			},
			Content:     item.ShortBody,
			PublishedAt: timestamppb.New(item.PublishedAt),
		}
	}
	return res
}
