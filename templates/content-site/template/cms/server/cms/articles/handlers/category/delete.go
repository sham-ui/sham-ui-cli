package category

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
	categoryRepository *repo.CategoryRepository
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
	err := h.categoryRepository.DeleteCategory(requestData.ID)
	if nil != err {
		return nil, fmt.Errorf("delete category: %s", err)
	}
	return &deleteResponse{
		Status: "Category deleted",
	}, nil
}

func NewDeleteHandler(db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	h := &deleteHandler{categoryRepository: repo.NewCategoryRepository(db)}
	return handler.Create(h, handler.WithOnlyForAuthenticated(sessionsStore))
}
