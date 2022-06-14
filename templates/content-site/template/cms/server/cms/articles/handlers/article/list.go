package article

import (
	"cms/core/handler"
	"cms/core/sessions"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type listHandler struct {
	db *sql.DB
}

type listHandlerData struct {
	offset int
	limit  int
}

func (h *listHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	offset, err := strconv.Atoi(ctx.Request.URL.Query().Get("offset"))
	if nil != err {
		offset = 0
	}
	limit, err := strconv.Atoi(ctx.Request.URL.Query().Get("limit"))
	if nil != err {
		limit = 20
	}
	return &listHandlerData{
		offset: offset,
		limit:  limit,
	}, nil
}

func (h *listHandler) Validate(_ *handler.Context, data interface{}) (*handler.Validation, error) {
	validation := handler.NewValidation()
	params := data.(*listHandlerData)
	if params.limit <= 0 {
		validation.AddError("limit must be > 0")
	}
	if params.offset < 0 {
		validation.AddError("offset must be >= 0")
	}
	return validation, nil
}

func (h *listHandler) getArticlesCount() (int, error) {
	var count int
	err := h.db.QueryRow("SELECT COUNT(*) FROM article").Scan(&count)
	if nil != err {
		return count, fmt.Errorf("select count: %s", err)
	}
	return count, nil
}

type articleListItemData struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	CategoryID  string    `json:"category_id"`
	PublishedAt time.Time `json:"published_at"`
}

func (h *listHandler) getArticles(offset, limit int) ([]*articleListItemData, error) {
	var articles []*articleListItemData
	rows, err := h.db.Query(
		"SELECT id, title, slug, category_id, published_at FROM article ORDER BY published_at DESC LIMIT $1 OFFSET $2",
		limit,
		offset,
	)
	defer rows.Close()
	if nil != err {
		return articles, fmt.Errorf("query: %s", err)
	}
	for rows.Next() {
		data := &articleListItemData{}
		err := rows.Scan(&data.ID, &data.Title, &data.Slug, &data.CategoryID, &data.PublishedAt)
		if nil != err {
			return articles, fmt.Errorf("scan row: %s", err)
		}
		articles = append(articles, data)
	}
	return articles, nil
}

func (h *listHandler) Process(_ *handler.Context, data interface{}) (interface{}, error) {
	params := data.(*listHandlerData)
	count, err := h.getArticlesCount()
	if nil != err {
		return nil, fmt.Errorf("articles count: %s", err)
	}
	articles, err := h.getArticles(params.offset, params.limit)
	if nil != err {
		return nil, fmt.Errorf("get articles: %s", err)
	}
	return map[string]interface{}{
		"meta": map[string]int{
			"offset": params.offset,
			"limit":  params.limit,
			"total":  count,
		},
		"articles": articles,
	}, nil
}

func NewListHandler(db *sql.DB, sessionStore *sessions.Store) http.HandlerFunc {
	return handler.Create(&listHandler{db: db}, handler.WithOnlyForAuthenticated(sessionStore))
}
