package app

import "cms/internal/controller/grpc"

type grpcServerDependencies struct {
	articleService    grpc.ArticleService
	categoryService   grpc.CategoryService
	tagService        grpc.TagService
	articleTagService grpc.ArticleTagService
	assetService      grpc.AssetService
}

func (r grpcServerDependencies) ArticleService() grpc.ArticleService { //nolint:ireturn
	return r.articleService
}

func (r grpcServerDependencies) CategoryService() grpc.CategoryService { //nolint:ireturn
	return r.categoryService
}

func (r grpcServerDependencies) TagService() grpc.TagService { //nolint:ireturn
	return r.tagService
}

func (r grpcServerDependencies) ArticleTagService() grpc.ArticleTagService { //nolint:ireturn
	return r.articleTagService
}

func (r grpcServerDependencies) AssetService() grpc.AssetService { //nolint:ireturn
	return r.assetService
}
