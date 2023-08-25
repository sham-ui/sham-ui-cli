package article

import (
	"cms/articles"
	repo "cms/articles/db"
	"cms/core/handler"
	"cms/core/sessions"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

type updateHandler struct {
	withArticleTags
	db                *sql.DB
	articleRepository *repo.ArticleRepository
}

func (h *updateHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	var data createUpdateRequestData
	err := json.NewDecoder(ctx.Request.Body).Decode(&data)
	if nil != err {
		return nil, fmt.Errorf("can't extract json data: %s", err)
	}
	data.Title = strings.TrimSpace(data.Title)
	data.Slug = articles.GenerateSlug(data.Title)
	id, err := strconv.Atoi(mux.Vars(ctx.Request)["id"])
	if nil != err {
		return nil, fmt.Errorf("can't get id")
	}
	data.ID = id
	return &data, nil
}

func (h *updateHandler) Validate(ctx *handler.Context, data interface{}) (*handler.Validation, error) {
	validation := handler.NewValidation()
	requestData := data.(*createUpdateRequestData)
	if "" == requestData.Title {
		validation.AddError("Title must not be empty.")
	}
	isUnique, id, err := h.articleRepository.IsUnique(requestData.Slug)
	if nil != err {
		return nil, fmt.Errorf("is unique slug: %s", err)
	}
	if !isUnique && id != requestData.ID {
		validation.AddError("Title is already in use.")
	}
	return validation, nil
}

func (h *updateHandler) Process(_ *handler.Context, data interface{}) (interface{}, error) {
	requestData := data.(*createUpdateRequestData)
	tx, err := h.db.Begin()
	if nil != err {
		return nil, fmt.Errorf("begin tx: %s", err)
	}
	tagIDs, err := h.getOrCreateTags(tx, requestData.Tags)
	if nil != err {
		tx.Rollback()
		return nil, fmt.Errorf("get or create tags: %s", err)
	}
	_, err = tx.Exec("UPDATE article SET "+
		"title = $2, slug = $3, category_id = $4, short_body = $5, body = $6, published_at = $7 "+
		"WHERE id = $1",
		requestData.ID,
		requestData.Title,
		requestData.Slug,
		requestData.CategoryID,
		requestData.ShortBody,
		requestData.Body,
		requestData.PublishedAt,
	)
	if nil != err {
		tx.Rollback()
		return nil, fmt.Errorf("insert article: %s", err)
	}

	existedTagIDs, err := h.articleRepository.GetTagIDs(tx, requestData.ID)
	if nil != err {
		tx.Rollback()
		return nil, fmt.Errorf("get existed tags: %s", err)
	}
	tagIDsToInsert, tagIDsToDelete := h.mergeTags(existedTagIDs, tagIDs)
	for _, tagID := range tagIDsToInsert {
		_, err := tx.Exec("INSERT INTO article_tag(article_id, tag_id) VALUES ($1, $2)", requestData.ID, tagID)
		if nil != err {
			tx.Rollback()
			return nil, fmt.Errorf("insert article_tag: %s", err)
		}
	}
	for _, tagID := range tagIDsToDelete {
		_, err := tx.Exec("DELETE FROM article_tag WHERE article_id = $1 AND tag_id = $2 ", requestData.ID, tagID)
		if nil != err {
			tx.Rollback()
			return nil, fmt.Errorf("delete article_tag: %s", err)
		}
	}

	err = tx.Commit()
	if nil != err {
		return nil, fmt.Errorf("commit tx: %s", err)
	}
	return &createUpdateResponse{
		Status: "Article updated",
	}, nil
}

func NewUpdateHandler(logger logr.Logger, db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	h := &updateHandler{
		withArticleTags:   withArticleTags{tagsRepository: repo.NewTagRepository(db)},
		db:                db,
		articleRepository: repo.NewArticleRepository(db),
	}
	return handler.Create(logger, h, handler.WithOnlyForAuthenticated(sessionsStore))
}
