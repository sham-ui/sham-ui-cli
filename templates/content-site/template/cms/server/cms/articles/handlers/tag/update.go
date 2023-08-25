package tag

import (
	"cms/articles"
	repo "cms/articles/db"
	"cms/core/handler"
	"cms/core/sessions"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type updateHandler struct {
	tagRepository *repo.TagRepository
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
	isUnique, err := h.tagRepository.IsUniqueTag(requestData.Slug)
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
	tagData := &repo.Tag{
		Name: requestData.Name,
		Slug: requestData.Slug,
	}
	err := h.tagRepository.UpdateTag(requestData.ID, tagData)
	if nil != err {
		return nil, fmt.Errorf("update tag: %s", err)
	}
	return &updateResponse{
		Status: "Tag updated",
	}, nil
}

func NewUpdateHandler(logger logr.Logger, db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	h := &updateHandler{tagRepository: repo.NewTagRepository(db)}
	return handler.Create(logger, h, handler.WithOnlyForAuthenticated(sessionsStore))
}
