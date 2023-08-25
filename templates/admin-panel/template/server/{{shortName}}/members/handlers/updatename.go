package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"net/http"
	"{{shortName}}/core/handler"
	"{{shortName}}/core/sessions"
	"{{shortName}}/members"
	"strings"
)

type updateNameHandler struct {
	db *sql.DB
}

type updateNameRequestData struct {
	NewName string
}

func (h *updateNameHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	var data updateNameRequestData
	err := json.NewDecoder(ctx.Request.Body).Decode(&data)
	if nil != err {
		return nil, fmt.Errorf("can't extract json data: %s", err)
	}
	data.NewName = strings.TrimSpace(data.NewName)
	return &data, nil
}

func (h *updateNameHandler) Validate(_ *handler.Context, data interface{}) (*handler.Validation, error) {
	validation := handler.NewValidation()
	requestData := data.(*updateNameRequestData)
	if "" == requestData.NewName {
		validation.AddError("Name must have more than 0 characters.")
	}
	return validation, nil
}

type updateNameResponse struct {
	Status string `json:"Status"`
}

func (h *updateNameHandler) Process(ctx *handler.Context, data interface{}) (interface{}, error) {
	requestData := data.(*updateNameRequestData)
	session, _ := ctx.GetSession()
	if err := members.UpdateMemberName(h.db, session.ID, requestData.NewName); nil != err {
		return nil, fmt.Errorf("can't update member name: %s", err)
	}
	rawSession := session.GetRawSession()
	rawSession.Values["name"] = requestData.NewName
	if err := rawSession.Save(ctx.Request, ctx.Response); err != nil {
		return nil, fmt.Errorf("can't save session: %s", err)
	}
	return &updateNameResponse{Status: "Name updated"}, nil
}

func NewUpdateNameHandler(logger logr.Logger, db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	return handler.Create(logger, &updateNameHandler{db: db}, handler.WithOnlyForAuthenticated(sessionsStore))
}
