package integration

import (
	"context"

	"site/internal/external_api/cms/proto"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name cmsServer --inpackage --testonly
type cmsServer interface {
	ArticleList(context.Context, *proto.ArticleListRequest) (*proto.ArticleListResponse, error)
	//nolint:lll
	ArticleListForCategory(context.Context, *proto.ArticleListForCategoryRequest) (*proto.ArticleListForCategoryResponse, error)
	ArticleListForTag(context.Context, *proto.ArticleListForTagRequest) (*proto.ArticleListForTagResponse, error)
	ArticleListForQuery(context.Context, *proto.ArticleListForQueryRequest) (*proto.ArticleListResponse, error)
	Article(context.Context, *proto.ArticleRequest) (*proto.ArticleResponse, error)
	Asset(context.Context, *proto.AssetRequest) (*proto.AssetResponse, error)
}
