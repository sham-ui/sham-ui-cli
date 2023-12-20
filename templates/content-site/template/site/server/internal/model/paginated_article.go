package model

type PaginatedArticles struct {
	Articles []ShortArticle
	Total    int64
}

type PaginatedArticleForTag struct {
	Articles []ShortArticle
	Total    int64
	TagName  string
}

type PaginatedArticleForCategory struct {
	Articles     []ShortArticle
	Total        int64
	CategoryName string
}
