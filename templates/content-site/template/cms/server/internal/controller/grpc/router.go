package grpc

import (
	articleDetailHandler "cms/internal/controller/grpc/handler/article/detail"
	articleListCategoryHandler "cms/internal/controller/grpc/handler/article/list/category"
	articleListDefaultHandler "cms/internal/controller/grpc/handler/article/list/detault"
	articleListQueryHandler "cms/internal/controller/grpc/handler/article/list/query"
	articleListTagHandler "cms/internal/controller/grpc/handler/article/list/tag"
	assetHandler "cms/internal/controller/grpc/handler/asset"
	"cms/internal/controller/grpc/proto"
	"context"
	"fmt"
)

type (
	ArticleListHandler interface {
		ArticleList(ctx context.Context, req *proto.ArticleListRequest) (*proto.ArticleListResponse, error)
	}
	ArticleListCategoryHandler interface {
		ArticleListForCategory(
			ctx context.Context,
			req *proto.ArticleListForCategoryRequest,
		) (*proto.ArticleListForCategoryResponse, error)
	}
	ArticleListTagHandler interface {
		ArticleListForTag(
			ctx context.Context,
			req *proto.ArticleListForTagRequest,
		) (*proto.ArticleListForTagResponse, error)
	}
	ArticleListQueryHandler interface {
		ArticleListForQuery(
			ctx context.Context,
			req *proto.ArticleListForQueryRequest,
		) (*proto.ArticleListResponse, error)
	}
	ArticleDetailHandler interface {
		Article(ctx context.Context, req *proto.ArticleRequest) (*proto.ArticleResponse, error)
	}
	AssetHandler interface {
		Asset(ctx context.Context, req *proto.AssetRequest) (*proto.AssetResponse, error)
	}
)

type router struct {
	proto.UnimplementedCMSServer
	articleListHandler         ArticleListHandler
	articleListCategoryHandler ArticleListCategoryHandler
	articleListTagHandler      ArticleListTagHandler
	articleListQueryHandler    ArticleListQueryHandler
	articleDetailHandler       ArticleDetailHandler
	assetHandler               AssetHandler
}

func (r *router) ArticleList(ctx context.Context, req *proto.ArticleListRequest) (*proto.ArticleListResponse, error) {
	resp, err := r.articleListHandler.ArticleList(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("article list: %w", err)
	}
	return resp, nil
}

func (r *router) ArticleListForCategory(
	ctx context.Context,
	req *proto.ArticleListForCategoryRequest,
) (*proto.ArticleListForCategoryResponse, error) {
	resp, err := r.articleListCategoryHandler.ArticleListForCategory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("article list for category: %w", err)
	}
	return resp, nil
}

func (r *router) ArticleListForTag(
	ctx context.Context,
	req *proto.ArticleListForTagRequest,
) (*proto.ArticleListForTagResponse, error) {
	resp, err := r.articleListTagHandler.ArticleListForTag(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("article list for tag: %w", err)
	}
	return resp, nil
}

func (r *router) ArticleListForQuery(
	ctx context.Context,
	req *proto.ArticleListForQueryRequest,
) (*proto.ArticleListResponse, error) {
	resp, err := r.articleListQueryHandler.ArticleListForQuery(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("article list for query: %w", err)
	}
	return resp, nil
}

func (r *router) Article(
	ctx context.Context,
	req *proto.ArticleRequest,
) (*proto.ArticleResponse, error) {
	resp, err := r.articleDetailHandler.Article(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("article: %w", err)
	}
	return resp, nil
}

func (r *router) Asset(
	ctx context.Context,
	req *proto.AssetRequest,
) (*proto.AssetResponse, error) {
	resp, err := r.assetHandler.Asset(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("asset: %w", err)
	}
	return resp, nil
}

func newRouter(deps HandlerDependencyProvider) *router {
	return &router{
		UnimplementedCMSServer: proto.UnimplementedCMSServer{},
		articleListHandler:     articleListDefaultHandler.New(deps.ArticleService()),
		articleListCategoryHandler: articleListCategoryHandler.New(
			deps.ArticleService(),
			deps.CategoryService(),
		),
		articleListTagHandler: articleListTagHandler.New(
			deps.ArticleService(),
			deps.TagService(),
		),
		articleListQueryHandler: articleListQueryHandler.New(deps.ArticleService()),
		articleDetailHandler: articleDetailHandler.NewHandler(
			deps.ArticleService(),
			deps.CategoryService(),
			deps.ArticleTagService(),
		),
		assetHandler: assetHandler.NewHandler(deps.AssetService()),
	}
}
