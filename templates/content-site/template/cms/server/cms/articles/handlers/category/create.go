package category

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
)

type createHandler struct {
	categoryRepository *repo.CategoryRepository
}

type createRequestData struct {
	Name string `json:"name"`
	Slug string `json:"-"`
}

func (h *createHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	var data createRequestData
	err := json.NewDecoder(ctx.Request.Body).Decode(&data)
	if nil != err {
		return nil, fmt.Errorf("can't extract json data: %s", err)
	}
	data.Name = strings.TrimSpace(data.Name)
	data.Slug = articles.GenerateSlug(data.Name)
	return &data, nil
}

func (h *createHandler) Validate(ctx *handler.Context, data interface{}) (*handler.Validation, error) {
	validation := handler.NewValidation()
	requestData := data.(*createRequestData)
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

type createResponse struct {
	Status string `json:"Status"`
}

func (h *createHandler) Process(_ *handler.Context, data interface{}) (interface{}, error) {
	requestData := data.(*createRequestData)
	categoryData := &repo.Category{
		Name: requestData.Name,
		Slug: requestData.Slug,
	}
	err := h.categoryRepository.CreateCategory(categoryData)
	if nil != err {
		return nil, fmt.Errorf("create category: %s", err)
	}
	return &createResponse{
		Status: "Category created",
	}, nil
}

func NewCreateHandler(db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	h := &createHandler{categoryRepository: repo.NewCategoryRepository(db)}
	return handler.Create(h, handler.WithOnlyForAuthenticated(sessionsStore))
}
