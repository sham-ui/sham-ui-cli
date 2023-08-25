package tag

import (
	"cms/core/handler"
	"cms/core/sessions"
	"database/sql"
	"fmt"
	"github.com/go-logr/logr"
	"net/http"
)

type listHandler struct {
	handler.HandlerWithoutExtractDataAndValidation
	db *sql.DB
}

type tagData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type listResponse struct {
	Categories []*tagData `json:"tags"`
}

func (h *listHandler) Process(_ *handler.Context, _ interface{}) (interface{}, error) {
	categories := make([]*tagData, 0)
	rows, err := h.db.Query("SELECT id, name, slug FROM tag ORDER BY id")
	defer rows.Close()
	if nil != err {
		return categories, fmt.Errorf("query: %s", err)
	}
	for rows.Next() {
		data := &tagData{}
		err := rows.Scan(&data.ID, &data.Name, &data.Slug)
		if nil != err {
			return categories, fmt.Errorf("scan row: %s", err)
		}
		categories = append(categories, data)
	}
	return &listResponse{
		Categories: categories,
	}, nil
}

func NewListHandler(logger logr.Logger, db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	return handler.Create(logger, &listHandler{db: db}, handler.WithOnlyForAuthenticated(sessionsStore))
}
