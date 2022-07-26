package cms

import (
	"context"
	"google.golang.org/grpc"
	"site/proto"
)

type MockResponse struct {
	Response interface{}
	Error    error
}

var _ proto.CMSClient = (*MockCMSClient)(nil)

type MockCMSClient struct {
	MockForArticleList            MockResponse
	MockForArticleListForCategory MockResponse
	MockForArticleListForTag      MockResponse
	MockForArticleListForQuery    MockResponse
	MockForArticle                MockResponse
	MockForAsset                  MockResponse
}

func (m *MockCMSClient) Asset(ctx context.Context, in *proto.AssetRequest, opts ...grpc.CallOption) (*proto.AssetResponse, error) {
	mock := m.MockForAsset
	return mock.Response.(*proto.AssetResponse), mock.Error
}

func (m *MockCMSClient) ArticleList(ctx context.Context, in *proto.ArticleListRequest, opts ...grpc.CallOption) (*proto.ArticleListResponse, error) {
	mock := m.MockForArticleList
	return mock.Response.(*proto.ArticleListResponse), mock.Error
}

func (m *MockCMSClient) ArticleListForCategory(ctx context.Context, in *proto.ArticleListForCategoryRequest, opts ...grpc.CallOption) (*proto.ArticleListForCategoryResponse, error) {
	mock := m.MockForArticleListForCategory
	return mock.Response.(*proto.ArticleListForCategoryResponse), mock.Error
}

func (m *MockCMSClient) ArticleListForTag(ctx context.Context, in *proto.ArticleListForTagRequest, opts ...grpc.CallOption) (*proto.ArticleListForTagResponse, error) {
	mock := m.MockForArticleListForTag
	return mock.Response.(*proto.ArticleListForTagResponse), mock.Error
}

func (m *MockCMSClient) ArticleListForQuery(ctx context.Context, in *proto.ArticleListForQueryRequest, opts ...grpc.CallOption) (*proto.ArticleListResponse, error) {
	mock := m.MockForArticleListForQuery
	return mock.Response.(*proto.ArticleListResponse), mock.Error
}

func (m *MockCMSClient) Article(ctx context.Context, in *proto.ArticleRequest, opts ...grpc.CallOption) (*proto.ArticleResponse, error) {
	mock := m.MockForArticle
	return mock.Response.(*proto.ArticleResponse), mock.Error
}
