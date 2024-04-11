package detail

import (
	"cms/internal/controller/grpc/proto"
	"cms/internal/model"
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type handler struct {
	articleService    ArticleService
	categoryService   CategoryService
	articleTagService ArticleTagService
}

func (h *handler) Article(ctx context.Context, request *proto.ArticleRequest) (*proto.ArticleResponse, error) {
	article, err := h.articleService.FindBySlug(ctx, model.ArticleSlug(request.Slug))
	switch {
	case errors.Is(err, model.ErrArticleNotFound):
		return &proto.ArticleResponse{
			Response: &proto.ArticleResponse_NotFound{
				NotFound: &proto.NotFound{},
			},
		}, nil
	case err != nil:
		return nil, fmt.Errorf("find article: %w", err)
	}
	category, err := h.categoryService.GetByID(ctx, article.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("category get by id: %w", err)
	}
	tags, err := h.articleTagService.GetTagForArticle(ctx, article.ID)
	if err != nil {
		return nil, fmt.Errorf("get tags for article: %w", err)
	}
	tagsResponse := make([]*proto.Tag, len(tags))
	for i := range tags {
		tagsResponse[i] = &proto.Tag{
			Name: tags[i].Name,
			Slug: string(tags[i].Slug),
		}
	}
	return &proto.ArticleResponse{
		Response: &proto.ArticleResponse_Article{
			Article: &proto.Article{
				Title:        article.Title,
				Slug:         string(article.Slug),
				Category:     &proto.Category{Name: category.Name, Slug: string(category.Slug)},
				ShortContent: article.ShortBody,
				Content:      article.Body,
				Tags:         tagsResponse,
				PublishedAt:  timestamppb.New(article.PublishedAt),
			},
		},
	}, nil
}

func NewHandler(
	articleService ArticleService,
	categoryService CategoryService,
	articleTagService ArticleTagService,
) *handler {
	return &handler{
		articleService:    articleService,
		categoryService:   categoryService,
		articleTagService: articleTagService,
	}
}
