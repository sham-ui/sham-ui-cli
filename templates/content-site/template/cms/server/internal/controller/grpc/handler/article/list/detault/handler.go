package detault

import (
	"cms/internal/controller/grpc/handler/article/list"
	"cms/internal/controller/grpc/proto"
	"context"
	"fmt"
)

type handler struct {
	articleService ArticleService
}

func (h *handler) ArticleList(
	ctx context.Context,
	request *proto.ArticleListRequest,
) (*proto.ArticleListResponse, error) {
	total, err := h.articleService.Total(ctx)
	if err != nil {
		return nil, fmt.Errorf("total: %w", err)
	}

	items, err := h.articleService.FindShortInfoWithCategory(ctx, request.Offset, request.Limit)
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}

	return &proto.ArticleListResponse{
		Articles: list.ArticleListItemFromModel(items),
		Total:    int64(total),
	}, nil
}

func New(articleService ArticleService) *handler {
	return &handler{
		articleService: articleService,
	}
}
