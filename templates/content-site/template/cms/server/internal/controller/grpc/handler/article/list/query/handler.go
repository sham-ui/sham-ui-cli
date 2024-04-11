package query

import (
	"cms/internal/controller/grpc/handler/article/list"
	"cms/internal/controller/grpc/proto"
	"context"
	"fmt"
)

type handler struct {
	articleService ArticleService
}

func (h *handler) ArticleListForQuery(
	ctx context.Context,
	request *proto.ArticleListForQueryRequest,
) (*proto.ArticleListResponse, error) {
	total, err := h.articleService.TotalForQuery(ctx, request.Query)
	if err != nil {
		return nil, fmt.Errorf("total: %w", err)
	}
	articles, err := h.articleService.FindShortInfoWithCategoryForQuery(
		ctx,
		request.Query,
		request.Offset,
		request.Limit,
	)
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}
	return &proto.ArticleListResponse{
		Total:    int64(total),
		Articles: list.ArticleListItemFromModel(articles),
	}, nil
}

func New(articleService ArticleService) *handler {
	return &handler{
		articleService: articleService,
	}
}
