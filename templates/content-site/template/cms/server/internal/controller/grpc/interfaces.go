package grpc

import (
	articleDetailHandler "cms/internal/controller/grpc/handler/article/detail"
	articleListCategoryHandler "cms/internal/controller/grpc/handler/article/list/category"
	articleListDefaultHandler "cms/internal/controller/grpc/handler/article/list/detault"
	articleListQueryHandler "cms/internal/controller/grpc/handler/article/list/query"
	articleListTagHandler "cms/internal/controller/grpc/handler/article/list/tag"
	assetHandler "cms/internal/controller/grpc/handler/asset"
)

type HandlerDependencyProvider interface {
	ArticleService() ArticleService
	CategoryService() CategoryService
	TagService() TagService
	ArticleTagService() ArticleTagService
	AssetService() AssetService
}

type ArticleService interface {
	articleListDefaultHandler.ArticleService
	articleListCategoryHandler.ArticleService
	articleListTagHandler.ArticleService
	articleListQueryHandler.ArticleService
	articleDetailHandler.ArticleService
}

type CategoryService interface {
	articleListCategoryHandler.CategoryService
	articleDetailHandler.CategoryService
}

type TagService interface {
	articleListTagHandler.TagService
}

type ArticleTagService interface {
	articleDetailHandler.ArticleTagService
}

type AssetService interface {
	assetHandler.AssetService
}
