package cms

import (
	"testing"
	"time"

	"site/internal/external_api/cms/proto"
	"site/internal/model"
	"site/pkg/asserts"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_articleListItemsToModel(t *testing.T) {
	// Arrange
	source := []*proto.ArticleListItem{
		{
			Title: "first-title",
			Slug:  "first-slug",
			Category: &proto.Category{
				Name: "first category",
				Slug: "first slug",
			},
			Content:     "first content",
			PublishedAt: timestamppb.New(time.Date(2022, time.April, 10, 13, 32, 10, 11, time.UTC)),
		},
		{
			Title: "second-title",
			Slug:  "second-slug",
			Category: &proto.Category{
				Name: "second category",
				Slug: "second slug",
			},
			Content:     "second content",
			PublishedAt: timestamppb.New(time.Date(2023, time.May, 11, 14, 33, 11, 12, time.UTC)),
		},
	}

	// Act
	result := articleListItemsToModel(source)

	// Assert
	asserts.Equals(t, []model.ShortArticle{
		{
			Title: "first-title",
			Slug:  "first-slug",
			Category: model.Category{
				Name: "first category",
				Slug: "first slug",
			},
			ShortContent: "first content",
			PublishedAt:  time.Date(2022, time.April, 10, 13, 32, 10, 11, time.UTC),
		},
		{
			Title: "second-title",
			Slug:  "second-slug",
			Category: model.Category{
				Name: "second category",
				Slug: "second slug",
			},
			ShortContent: "second content",
			PublishedAt:  time.Date(2023, time.May, 11, 14, 33, 11, 12, time.UTC),
		},
	}, result)
}
