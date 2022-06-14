package category

import (
	"cms/articles"
	repo "cms/articles/db"
	"cms/core/handler"
	"cms/core/sessions"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type updateHandler struct {
	categoryRepository *repo.CategoryRepository
}

type updateRequestData struct {
	ID   string `json:"-"`
	Name string `json:"name"`
	Slug string `json:"-"`
}

func (h *updateHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	var data updateRequestData
	err := json.NewDecoder(ctx.Request.Body).Decode(&data)
	if nil != err {
		return nil, fmt.Errorf("can't extract json data: %s", err)
	}
	data.Name = strings.TrimSpace(data.Name)
	data.Slug = articles.GenerateSlug(data.Name)
	data.ID = mux.Vars(ctx.Request)["id"]
	return &data, nil
}

func (h *updateHandler) Validate(ctx *handler.Context, data interface{}) (*handler.Validation, error) {
	validation := handler.NewValidation()
	requestData := data.(*updateRequestData)
	if "" == requestData.Name {
		validation.AddError("Name must not be empty.")
	}
	isUnique, err := h.categoryRepository.IsUniqueCategory(requestData.Slug)
	if nil != err {
		return nil, fmt.Errorf("is unique slug: %s", err)
	}
	if !isUnique {
		validation.AddError("Name is already in use.")
	}
	return validation, nil
}

type updateResponse struct {
	Status string `json:"Status"`
}

func (h *updateHandler) Process(_ *handler.Context, data interface{}) (interface{}, error) {
	requestData := data.(*updateRequestData)
	categoryData := &repo.Category{
		Name: requestData.Name,
		Slug: requestData.Slug,
	}
	err := h.categoryRepository.UpdateCategory(requestData.ID, categoryData)
	if nil != err {
		return nil, fmt.Errorf("update category: %s", err)
	}
	return &updateResponse{
		Status: "Category updated",
	}, nil
}

func NewUpdateHandler(db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	h := &updateHandler{categoryRepository: repo.NewCategoryRepository(db)}
	return handler.Create(h, handler.WithOnlyForAuthenticated(sessionsStore))
}
