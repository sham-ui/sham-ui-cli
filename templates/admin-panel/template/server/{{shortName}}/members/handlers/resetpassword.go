package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"net/http"
	"{{shortName}}/core/handler"
	"{{shortName}}/core/sessions"
	"{{shortName}}/members"
)

type resetPasswordHandler struct {
	db *sql.DB
}

type resetPasswordBody struct {
	Password1 string `json:"pass1"`
	Password2 string `json:"pass2"`
}

type resetPasswordRequestData struct {
	id   string
	body *resetPasswordBody
}

func (h *resetPasswordHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	var body resetPasswordBody
	err := json.NewDecoder(ctx.Request.Body).Decode(&body)
	if nil != err {
		return nil, fmt.Errorf("can't extract json data: %s", err)
	}
	return &resetPasswordRequestData{
		id:   mux.Vars(ctx.Request)["id"],
		body: &body,
	}, nil
}

func (h *resetPasswordHandler) Validate(_ *handler.Context, data interface{}) (*handler.Validation, error) {
	requestData := data.(*resetPasswordRequestData)
	validation := handler.NewValidation()
	if "" == requestData.body.Password1 || "" == requestData.body.Password2 {
		validation.AddError("Password must have more than 0 characters.")
	}
	if requestData.body.Password1 != requestData.body.Password2 {
		validation.AddError("Passwords don't match.")
	}
	return validation, nil
}

type resetPasswordResponse struct {
	Status string `json:"Status"`
}

func (h *resetPasswordHandler) Process(ctx *handler.Context, data interface{}) (interface{}, error) {
	requestData := data.(*resetPasswordRequestData)
	hashedPassword, err := members.HashPassword(requestData.body.Password1)
	if nil != err {
		return nil, fmt.Errorf("can't hash password: %s", err)
	}
	if err = members.UpdateMemberPassword(h.db, requestData.id, hashedPassword); nil != err {
		return nil, fmt.Errorf("can't update member password: %s", err)
	}
	return &resetPasswordResponse{Status: "Password updated"}, nil
}

func NewResetPasswordHandler(logger logr.Logger, db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	return handler.Create(logger, &resetPasswordHandler{db: db}, handler.WithOnlyForSuperuser(sessionsStore))
}
