package handlers

import (
	"database/sql"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"net/http"
	"{{shortName}}/core/handler"
	"{{shortName}}/core/sessions"
	"{{shortName}}/members"
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
	requestData := data.(*deleteRequestData)
	validation := handler.NewValidation()
	exists, err := members.HasMemberForID(h.db, requestData.ID)
	if nil != err {
		return nil, fmt.Errorf("has member for id: %s", err)
	}
	if !exists {
		validation.AddError("Member not exists.")
	}
	return validation, nil
}

type deleteResponse struct {
	Status string `json:"Status"`
}

func (h *deleteHandler) Process(_ *handler.Context, data interface{}) (interface{}, error) {
	requestData := data.(*deleteRequestData)
	err := members.DeleteMember(h.db, requestData.ID)
	if nil != err {
		return nil, fmt.Errorf("delete member: %s", err)
	}
	return &deleteResponse{
		Status: "Member deleted",
	}, nil
}

func NewDeleteHandler(logger logr.Logger, db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	return handler.Create(logger, &deleteHandler{db: db}, handler.WithOnlyForSuperuser(sessionsStore))
}
