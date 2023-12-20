package integration

import (
	"context"

	"site/internal/external_api/cms/proto"
)

type testCmsServer struct {
	proto.UnimplementedCMSServer
	m cmsServer
}

func (t testCmsServer) ArticleList(ctx context.Context, request *proto.ArticleListRequest) (*proto.ArticleListResponse, error) {
	return t.m.ArticleList(ctx, request)
}

func (t testCmsServer) ArticleListForCategory(ctx context.Context, request *proto.ArticleListForCategoryRequest) (*proto.ArticleListForCategoryResponse, error) {
	return t.m.ArticleListForCategory(ctx, request)
}

func (t testCmsServer) ArticleListForTag(ctx context.Context, request *proto.ArticleListForTagRequest) (*proto.ArticleListForTagResponse, error) {
	return t.m.ArticleListForTag(ctx, request)
}

func (t testCmsServer) ArticleListForQuery(ctx context.Context, request *proto.ArticleListForQueryRequest) (*proto.ArticleListResponse, error) {
	return t.m.ArticleListForQuery(ctx, request)
}

func (t testCmsServer) Article(ctx context.Context, request *proto.ArticleRequest) (*proto.ArticleResponse, error) {
	return t.m.Article(ctx, request)
}

func (t testCmsServer) Asset(ctx context.Context, request *proto.AssetRequest) (*proto.AssetResponse, error) {
	return t.m.Asset(ctx, request)
}
