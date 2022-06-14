package tag

import (
	repo "cms/articles/db"
	"cms/core/handler"
	"cms/core/sessions"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type deleteHandler struct {
	tagRepository *repo.TagRepository
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
	err := h.tagRepository.DeleteTag(requestData.ID)
	if nil != err {
		return nil, fmt.Errorf("delete tag: %s", err)
	}
	return &deleteResponse{
		Status: "Tag deleted",
	}, nil
}

func NewDeleteHandler(db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	h := &deleteHandler{tagRepository: repo.NewTagRepository(db)}
	return handler.Create(h, handler.WithOnlyForAuthenticated(sessionsStore))
}
