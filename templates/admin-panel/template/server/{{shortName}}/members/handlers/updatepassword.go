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
)

type updatePasswordHandler struct {
	db *sql.DB
}

type updatePasswordRequestData struct {
	NewPassword1, NewPassword2 string
}

func (h *updatePasswordHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	var data updatePasswordRequestData
	err := json.NewDecoder(ctx.Request.Body).Decode(&data)
	if nil != err {
		return nil, fmt.Errorf("can't extract json data: %s", err)
	}
	return &data, nil
}

func (h *updatePasswordHandler) Validate(_ *handler.Context, data interface{}) (*handler.Validation, error) {
	requestData := data.(*updatePasswordRequestData)
	validation := handler.NewValidation()
	if "" == requestData.NewPassword1 || "" == requestData.NewPassword2 {
		validation.AddError("Password must have more than 0 characters.")
	}
	if requestData.NewPassword1 != requestData.NewPassword2 {
		validation.AddError("Passwords don't match.")
	}
	return validation, nil
}

type updatePasswordResponse struct {
	Status string `json:"Status"`
}

func (h *updatePasswordHandler) Process(ctx *handler.Context, data interface{}) (interface{}, error) {
	requestData := data.(*updatePasswordRequestData)
	hashedPassword, err := members.HashPassword(requestData.NewPassword1)
	if nil != err {
		return nil, fmt.Errorf("can't hash password: %s", err)
	}
	session, _ := ctx.GetSession()
	if err = members.UpdateMemberPassword(h.db, session.ID, hashedPassword); nil != err {
		return nil, fmt.Errorf("can't update member password: %s", err)
	}
	return &updatePasswordResponse{Status: "Password updated"}, nil
}

func NewUpdatePasswordHandler(logger logr.Logger, db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	return handler.Create(logger, &updatePasswordHandler{db: db}, handler.WithOnlyForAuthenticated(sessionsStore))
}
