package http

import (
	"site/internal/controller/http/handler/api/articles/detail"
	"site/internal/controller/http/handler/api/articles/list/category"
	default_ "site/internal/controller/http/handler/api/articles/list/default"
	"site/internal/controller/http/handler/api/articles/list/query"
	"site/internal/controller/http/handler/api/articles/list/tag"
	"site/internal/controller/http/handler/assets"
	"site/internal/controller/http/handler/ssr"
)

type AssetsService interface {
	assets.Service
}

type ArticlesService interface {
	detail.ArticlesService
	default_.ArticlesService
	query.ArticlesService
	category.ArticlesService
	tag.ArticlesService
}

type ServerSideRender interface {
	ssr.ServerSideRender
}
