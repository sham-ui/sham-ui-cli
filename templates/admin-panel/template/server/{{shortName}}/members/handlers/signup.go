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

type signupHandler struct {
	db            *sql.DB
	sessionsStore *sessions.Store
}

type signupRequestData struct {
	Name, Email, Password, Password2 string
}

func (h *signupHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	var data signupRequestData
	err := json.NewDecoder(ctx.Request.Body).Decode(&data)
	if nil != err {
		return nil, fmt.Errorf("can't extract json data: %s", err)
	}
	data.Name = strings.TrimSpace(data.Name)
	data.Email = strings.TrimSpace(data.Email)
	return &data, nil
}

func (h *signupHandler) Validate(ctx *handler.Context, data interface{}) (*handler.Validation, error) {
	validation := handler.NewValidation()
	session, err := h.sessionsStore.GetSession(ctx.Request)
	if nil != err {
		return nil, fmt.Errorf("can't get session: %s", err)
	}
	if nil != session {
		validation.AddError("Already logged")
	}
	requestData := data.(*signupRequestData)
	if "" == requestData.Name {
		validation.AddError("Name must not be empty.")
	}
	if "" == requestData.Email {
		validation.AddError("Email must not be empty.")
	}
	if requestData.Password != requestData.Password2 {
		validation.AddError("Passwords do not match.")
	}
	isUnique, err := members.IsUniqueEmail(h.db, requestData.Email)
	if nil != err {
		return nil, fmt.Errorf("is unique email: %s", err)
	}
	if !isUnique {
		validation.AddError("Email is already in use.")
	}
	return validation, nil
}

type signupResponse struct {
	Status string `json:"Status"`
}

func (h *signupHandler) Process(_ *handler.Context, data interface{}) (interface{}, error) {
	requestData := data.(*signupRequestData)
	hashedPw, err := members.HashPassword(requestData.Password)
	if nil != err {
		return nil, fmt.Errorf("can't hash password: %s", err)
	}
	memberData := &members.MemberData{
		Name:        requestData.Name,
		Email:       requestData.Email,
		Password:    hashedPw,
		IsSuperuser: false,
	}
	err = members.CreateMember(h.db, memberData)
	if nil != err {
		return nil, fmt.Errorf("create member: %s", err)
	}
	return &signupResponse{
		Status: "Member created",
	}, nil
}

func NewSignupHandler(logger logr.Logger, db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	return handler.Create(logger, &signupHandler{db: db, sessionsStore: sessionsStore})
}
