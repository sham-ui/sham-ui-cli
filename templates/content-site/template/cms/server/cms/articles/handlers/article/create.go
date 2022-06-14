package article

import (
	"cms/articles"
	repo "cms/articles/db"
	"cms/core/handler"
	"cms/core/sessions"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type createUpdateRequestData struct {
	ID          int          `json:"-"`
	Title       string       `json:"title"`
	Slug        string       `json:"-"`
	CategoryID  string       `json:"category_id"`
	ShortBody   string       `json:"short_body"`
	Body        string       `json:"body"`
	PublishedAt time.Time    `json:"published_at"`
	Tags        []articleTag `json:"tags"`
}

type createUpdateResponse struct {
	Status string `json:"Status"`
}

type createHandler struct {
	withArticleTags
	db                *sql.DB
	articleRepository *repo.ArticleRepository
}

func (h *createHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	var data createUpdateRequestData
	err := json.NewDecoder(ctx.Request.Body).Decode(&data)
	if nil != err {
		return nil, fmt.Errorf("can't extract json data: %s", err)
	}
	data.Title = strings.TrimSpace(data.Title)
	data.Slug = articles.GenerateSlug(data.Title)
	return &data, nil
}

func (h *createHandler) Validate(ctx *handler.Context, data interface{}) (*handler.Validation, error) {
	validation := handler.NewValidation()
	requestData := data.(*createUpdateRequestData)
	if "" == requestData.Title {
		validation.AddError("Title must not be empty.")
	}
	isUnique, _, err := h.articleRepository.IsUnique(requestData.Slug)
	if nil != err {
		return nil, fmt.Errorf("is unique slug: %s", err)
	}
	if !isUnique {
		validation.AddError("Title is already in use.")
	}
	return validation, nil
}

func (h *createHandler) Process(_ *handler.Context, data interface{}) (interface{}, error) {
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
	var articleID int64
	row := tx.QueryRow("INSERT INTO "+
		"article(title, slug, category_id, short_body, body, published_at) "+
		"VALUES ($1,$2, $3, $4, $5, $6) RETURNING id",
		requestData.Title,
		requestData.Slug,
		requestData.CategoryID,
		requestData.ShortBody,
		requestData.Body,
		requestData.PublishedAt,
	)
	err = row.Scan(&articleID)
	if nil != err {
		tx.Rollback()
		return nil, fmt.Errorf("insert article: %s", err)
	}

	for _, tagID := range tagIDs {
		_, err := tx.Exec("INSERT INTO article_tag(article_id, tag_id) VALUES ($1, $2)", articleID, tagID)
		if nil != err {
			tx.Rollback()
			return nil, fmt.Errorf("insert article_tag: %s", err)
		}
	}

	err = tx.Commit()
	if nil != err {
		return nil, fmt.Errorf("commit tx: %s", err)
	}
	return &createUpdateResponse{
		Status: "Article created",
	}, nil
}

func NewCreateHandler(db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	h := &createHandler{
		withArticleTags:   withArticleTags{tagsRepository: repo.NewTagRepository(db)},
		db:                db,
		articleRepository: repo.NewArticleRepository(db),
	}
	return handler.Create(h, handler.WithOnlyForAuthenticated(sessionsStore))
}
