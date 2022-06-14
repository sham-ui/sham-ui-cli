package category

import (
	"cms/core/handler"
	"cms/core/sessions"
	"database/sql"
	"fmt"
	"net/http"
)

type listHandler struct {
	handler.HandlerWithoutExtractDataAndValidation
	db *sql.DB
}

type categoryData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type listResponse struct {
	Categories []*categoryData `json:"categories"`
}

func (h *listHandler) Process(_ *handler.Context, _ interface{}) (interface{}, error) {
	categories := make([]*categoryData, 0)
	rows, err := h.db.Query("SELECT id, name, slug FROM category ORDER BY id")
	defer rows.Close()
	if nil != err {
		return categories, fmt.Errorf("query: %s", err)
	}
	for rows.Next() {
		data := &categoryData{}
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

func NewListHandler(db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	return handler.Create(&listHandler{db: db}, handler.WithOnlyForAuthenticated(sessionsStore))
}
