package article

import (
	"cms/core/handler"
	"cms/core/sessions"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type detailResponse struct {
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	CategoryID  string    `json:"category_id"`
	ShortBody   string    `json:"short_body"`
	Body        string    `json:"body"`
	PublishedAt time.Time `json:"published_at"`
	Tags        []string  `json:"tags"`
}

type detailHandler struct {
	db *sql.DB
}

func (h *detailHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	id, err := strconv.Atoi(mux.Vars(ctx.Request)["id"])
	if nil != err {
		return nil, fmt.Errorf("can't get id")
	}
	return id, nil
}

func (h *detailHandler) Validate(ctx *handler.Context, data interface{}) (*handler.Validation, error) {
	return handler.NewValidation(), nil
}

func (h *detailHandler) Process(_ *handler.Context, data interface{}) (interface{}, error) {
	id := data.(int)
	article := &detailResponse{}
	err := h.db.
		QueryRow(
			"SELECT title, slug, category_id, short_body, body, published_at FROM article WHERE id = $1",
			id,
		).
		Scan(&article.Title, &article.Slug, &article.CategoryID, &article.ShortBody, &article.Body, &article.PublishedAt)
	if nil != err {
		return nil, fmt.Errorf("query article: %s", err)
	}
	rows, err := h.db.Query(
		"SELECT t.slug FROM article_tag at INNER JOIN tag t ON t.id = at.tag_id WHERE at.article_id = $1",
		id,
	)
	if nil != err {
		return nil, fmt.Errorf("query tags: %s", err)
	}
	for rows.Next() {
		var tagSlug string
		err := rows.Scan(&tagSlug)
		if nil != err {
			return nil, fmt.Errorf("scan tag: %s", err)
		}
		article.Tags = append(article.Tags, tagSlug)
	}
	return article, nil
}

func NewDetailHandler(db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	h := &detailHandler{
		db: db,
	}
	return handler.Create(h, handler.WithOnlyForAuthenticated(sessionsStore))
}
