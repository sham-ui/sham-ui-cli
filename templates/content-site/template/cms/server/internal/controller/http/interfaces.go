package http

import (
	membersCreate "cms/internal/controller/http/handler/api/admin/members/create"
	membersList "cms/internal/controller/http/handler/api/admin/members/list"
	membersDelete "cms/internal/controller/http/handler/api/admin/members/remove"
	membersResetPassword "cms/internal/controller/http/handler/api/admin/members/reset_password"
	membersUpdate "cms/internal/controller/http/handler/api/admin/members/update"
	articlesCreate "cms/internal/controller/http/handler/api/articles/articles/create"
	articlesDetail "cms/internal/controller/http/handler/api/articles/articles/detail"
	articlesList "cms/internal/controller/http/handler/api/articles/articles/list"
	articlesDelete "cms/internal/controller/http/handler/api/articles/articles/remove"
	articlesUpdate "cms/internal/controller/http/handler/api/articles/articles/update"
	articlesCategoryCreate "cms/internal/controller/http/handler/api/articles/categories/create"
	articlesCategoryList "cms/internal/controller/http/handler/api/articles/categories/list"
	articlesCategoryRemove "cms/internal/controller/http/handler/api/articles/categories/remove"
	articlesCategoryUpdate "cms/internal/controller/http/handler/api/articles/categories/update"
	articlesTagCreate "cms/internal/controller/http/handler/api/articles/tags/create"
	articlesTagList "cms/internal/controller/http/handler/api/articles/tags/list"
	articlesTagRemove "cms/internal/controller/http/handler/api/articles/tags/remove"
	articlesTagUpdate "cms/internal/controller/http/handler/api/articles/tags/update"
	assetsUpload "cms/internal/controller/http/handler/api/assets/upload"
	"cms/internal/controller/http/handler/api/members/email"
	"cms/internal/controller/http/handler/api/members/name"
	"cms/internal/controller/http/handler/api/members/password"
	"cms/internal/controller/http/handler/api/session/login"
	"cms/internal/controller/http/handler/api/session/logout"
	sessionMW "cms/internal/controller/http/middleware/session"
)

//nolint:lll
//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name HandlerDependencyProvider --inpackage --testonly --with-expecter
type HandlerDependencyProvider interface {
	PasswordService() PasswordService
	SessionService() SessionService
	SlugifyService() SlugifyService
	MemberService() MemberService
	ArticleCategoryService() ArticleCategoryService
	ArticleTagService() ArticleTagService
	ArticleService() ArticleService
	AssetsService() AssetsService
}

type PasswordService interface {
	login.PasswordService
	password.PasswordService
	membersCreate.PasswordService
	membersResetPassword.PasswordService
}

type SessionService interface {
	sessionMW.Service
	login.SessionService
	logout.Service
	email.SessionService
	name.SessionService
}

type SlugifyService interface {
	articlesCategoryCreate.SlugifyService
	articlesCategoryUpdate.SlugifyService
	articlesTagCreate.SlugifyService
	articlesTagUpdate.SlugifyService
	articlesCreate.SlugifyService
	articlesUpdate.SlugifyService
}

type MemberService interface {
	login.MemberService
	email.MemberService
	name.MemberService
	password.MemberService
	membersList.MemberService
	membersCreate.MemberService
	membersUpdate.MemberService
	membersResetPassword.MemberService
	membersDelete.MemberService
}

type ArticleCategoryService interface {
	articlesCategoryCreate.CategoryService
	articlesCategoryUpdate.CategoryService
	articlesCategoryRemove.CategoryService
	articlesCategoryList.CategoryService
}

type ArticleTagService interface {
	articlesTagCreate.TagService
	articlesTagUpdate.TagService
	articlesTagRemove.TagService
	articlesTagList.TagService
}

type ArticleService interface {
	articlesCreate.ArticleService
	articlesUpdate.ArticleService
	articlesList.ArticleService
	articlesDetail.ArticleService
	articlesDelete.ArticleService
}

type AssetsService interface {
	assetsUpload.AssetsService
}
