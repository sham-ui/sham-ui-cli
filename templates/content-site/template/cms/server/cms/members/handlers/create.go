package handlers

import (
	"cms/core/handler"
	"cms/core/sessions"
	"cms/members"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"net/http"
	"strings"
)

type createHandler struct {
	db *sql.DB
}

type createRequestData struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsSuperUser bool   `json:"is_superuser"`
}

func (h *createHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	var data createRequestData
	err := json.NewDecoder(ctx.Request.Body).Decode(&data)
	if nil != err {
		return nil, fmt.Errorf("can't extract json data: %s", err)
	}
	data.Name = strings.TrimSpace(data.Name)
	data.Email = strings.TrimSpace(data.Email)
	return &data, nil
}

func (h *createHandler) Validate(ctx *handler.Context, data interface{}) (*handler.Validation, error) {
	validation := handler.NewValidation()
	requestData := data.(*createRequestData)
	if "" == requestData.Name {
		validation.AddError("Name must not be empty.")
	}
	if "" == requestData.Email {
		validation.AddError("Email must not be empty")
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

type createResponse struct {
	Status string `json:"Status"`
}

func (h *createHandler) Process(_ *handler.Context, data interface{}) (interface{}, error) {
	requestData := data.(*createRequestData)
	hashedPw, err := members.HashPassword(requestData.Password)
	if nil != err {
		return nil, fmt.Errorf("can't hash password: %s", err)
	}
	memberData := &members.MemberData{
		Name:        requestData.Name,
		Email:       requestData.Email,
		Password:    hashedPw,
		IsSuperuser: requestData.IsSuperUser,
	}
	err = members.CreateMember(h.db, memberData)
	if nil != err {
		return nil, fmt.Errorf("create member: %s", err)
	}
	return &createResponse{
		Status: "Member created",
	}, nil
}

func NewCreateHandler(logger logr.Logger, db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	return handler.Create(logger, &createHandler{db: db}, handler.WithOnlyForSuperuser(sessionsStore))
}
