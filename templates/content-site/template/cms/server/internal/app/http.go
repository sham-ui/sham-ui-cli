package app

import "cms/internal/controller/http"

type httpServerDependencies struct {
	passwordService        http.PasswordService
	sessionService         http.SessionService
	slugifyService         http.SlugifyService
	memberService          http.MemberService
	articleCategoryService http.ArticleCategoryService
	articleTagService      http.ArticleTagService
	articleService         http.ArticleService
	assetsService          http.AssetsService
}

func (r httpServerDependencies) PasswordService() http.PasswordService { //nolint:ireturn
	return r.passwordService
}

func (r httpServerDependencies) SessionService() http.SessionService { //nolint:ireturn
	return r.sessionService
}

func (r httpServerDependencies) SlugifyService() http.SlugifyService { //nolint:ireturn
	return r.slugifyService
}

func (r httpServerDependencies) MemberService() http.MemberService { //nolint:ireturn
	return r.memberService
}

func (r httpServerDependencies) ArticleCategoryService() http.ArticleCategoryService { //nolint:ireturn
	return r.articleCategoryService
}

func (r httpServerDependencies) ArticleTagService() http.ArticleTagService { //nolint:ireturn
	return r.articleTagService
}

func (r httpServerDependencies) ArticleService() http.ArticleService { //nolint:ireturn
	return r.articleService
}

func (r httpServerDependencies) AssetsService() http.AssetsService { //nolint:ireturn
	return r.assetsService
}
