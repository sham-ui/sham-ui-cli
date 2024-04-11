package tag

import (
	"cms/internal/controller/grpc/handler/article/list"
	"cms/internal/controller/grpc/proto"
	"cms/internal/model"
	"context"
	"fmt"
)

type handler struct {
	articleService ArticleService
	tagService     TagService
}

func (h *handler) ArticleListForTag(
	ctx context.Context,
	request *proto.ArticleListForTagRequest,
) (*proto.ArticleListForTagResponse, error) {
	slug := model.TagSlug(request.TagSlug)

	tag, err := h.tagService.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("tag get by slug: %w", err)
	}

	total, err := h.articleService.TotalForTag(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("total: %w", err)
	}

	items, err := h.articleService.FindShortInfoWithCategoryForTag(ctx, slug, request.Offset, request.Limit)
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}

	return &proto.ArticleListForTagResponse{
		Articles: list.ArticleListItemFromModel(items),
		Total:    int64(total),
		TagName:  tag.Name,
	}, nil
}

func New(
	articleService ArticleService,
	tagService TagService,
) *handler {
	return &handler{
		articleService: articleService,
		tagService:     tagService,
	}
}
