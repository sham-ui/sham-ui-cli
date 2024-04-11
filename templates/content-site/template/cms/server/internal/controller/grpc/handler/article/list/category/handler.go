package category

import (
	"cms/internal/controller/grpc/handler/article/list"
	"cms/internal/controller/grpc/proto"
	"cms/internal/model"
	"context"
	"fmt"
)

type handler struct {
	articleService  ArticleService
	categoryService CategoryService
}

func (h *handler) ArticleListForCategory(
	ctx context.Context,
	request *proto.ArticleListForCategoryRequest,
) (*proto.ArticleListForCategoryResponse, error) {
	slug := model.CategorySlug(request.CategorySlug)

	cat, err := h.categoryService.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("category get by slug: %w", err)
	}

	total, err := h.articleService.TotalForCategory(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("total: %w", err)
	}

	items, err := h.articleService.FindShortInfoWithCategoryForCategory(ctx, slug, request.Offset, request.Limit)
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}

	return &proto.ArticleListForCategoryResponse{
		Articles:     list.ArticleListItemFromModel(items),
		Total:        int64(total),
		CategoryName: cat.Name,
	}, nil
}

func New(articleService ArticleService, categoryService CategoryService) *handler {
	return &handler{
		articleService:  articleService,
		categoryService: categoryService,
	}
}
