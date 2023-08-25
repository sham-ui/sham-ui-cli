package article

import (
	"cms/core/handler"
	"cms/core/sessions"
	"database/sql"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"net/http"
)

type deleteHandler struct {
	db *sql.DB
}

type deleteRequestData struct {
	ID string
}

func (h *deleteHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	return &deleteRequestData{
		ID: mux.Vars(ctx.Request)["id"],
	}, nil
}

func (h *deleteHandler) Validate(ctx *handler.Context, data interface{}) (*handler.Validation, error) {
	return handler.NewValidation(), nil
}

type deleteResponse struct {
	Status string `json:"Status"`
}

func (h *deleteHandler) Process(_ *handler.Context, data interface{}) (interface{}, error) {
	requestData := data.(*deleteRequestData)
	_, err := h.db.Exec("DELETE FROM article WHERE id = $1", requestData.ID)
	if nil != err {
		return nil, fmt.Errorf("delete article: %s", err)
	}
	return &deleteResponse{
		Status: "Article deleted",
	}, nil
}

func NewDeleteHandler(logger logr.Logger, db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	h := &deleteHandler{db: db}
	return handler.Create(logger, h, handler.WithOnlyForAuthenticated(sessionsStore))
}
