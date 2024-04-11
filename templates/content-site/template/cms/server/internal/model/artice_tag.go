package model

import "errors"

type ArticleTag struct {
	ArticleID ArticleID
	TagID     TagID
}

var ErrArticleTagAlreadyExists = errors.New("article tag already exists")
