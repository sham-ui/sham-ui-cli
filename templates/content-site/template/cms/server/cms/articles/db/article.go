package db

import (
	"database/sql"
	"fmt"
)

type ArticleRepository struct {
	db *sql.DB
}

func (r *ArticleRepository) IsUnique(slug string) (bool, int, error) {
	var existingId int
	row := r.db.QueryRow("SELECT id FROM article WHERE slug = $1", slug)
	err := row.Scan(&existingId)
	if err == sql.ErrNoRows {
		return true, 0, nil
	}
	if err != nil {
		return false, 0, fmt.Errorf("select id: %s", err)
	}
	return false, existingId, nil
}

func (r *ArticleRepository) GetTagIDs(tx *sql.Tx, articleID int) ([]int, error) {
	res, err := tx.Query("SELECT tag_id FROM article_tag WHERE article_id = $1", articleID)
	if nil != err {
		return nil, fmt.Errorf("select tag_id: %s", err)
	}
	var tagIDs []int
	for res.Next() {
		var tagID int
		err := res.Scan(&tagID)
		if nil != err {
			return nil, fmt.Errorf("scan tag_id: %s", err)
		}
		tagIDs = append(tagIDs, tagID)
	}
	return tagIDs, nil
}

func NewArticleRepository(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}
