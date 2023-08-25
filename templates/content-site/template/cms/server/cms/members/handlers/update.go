package handlers

import (
	"cms/core/handler"
	"cms/core/sessions"
	"cms/members"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type updateHandler struct {
	db *sql.DB
}

type updateRequestDataBody struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	IsSuperUser bool   `json:"is_superuser"`
}

type updateRequestData struct {
	id   string
	body *updateRequestDataBody
}

func (h *updateHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	var body updateRequestDataBody
	err := json.NewDecoder(ctx.Request.Body).Decode(&body)
	if nil != err {
		return nil, fmt.Errorf("can't extract json data: %s", err)
	}
	body.Name = strings.TrimSpace(body.Name)
	body.Email = strings.TrimSpace(body.Email)
	return &updateRequestData{
		id:   mux.Vars(ctx.Request)["id"],
		body: &body,
	}, nil
}

func (h *updateHandler) Validate(ctx *handler.Context, data interface{}) (*handler.Validation, error) {
	validation := handler.NewValidation()
	requestData := data.(*updateRequestData)
	if "" == requestData.body.Name {
		validation.AddError("Name must not be empty.")
	}
	if "" == requestData.body.Email {
		validation.AddError("Email must not be empty")
	}
	isUnique, err := members.IsUniqueEmailForMemberID(h.db, requestData.id, requestData.body.Email)
	if nil != err {
		return nil, fmt.Errorf("is unique email: %s", err)
	}
	if !isUnique {
		validation.AddError("Email is already in use.")
	}
	return validation, nil
}

type updateResponse struct {
	Status string `json:"Status"`
}

func (h *updateHandler) Process(_ *handler.Context, data interface{}) (interface{}, error) {
	requestData := data.(*updateRequestData)
	err := members.UpdateMemberData(h.db, &members.MemberData{
		ID:          requestData.id,
		Name:        requestData.body.Name,
		Email:       requestData.body.Email,
		IsSuperuser: requestData.body.IsSuperUser,
	})
	if nil != err {
		return nil, fmt.Errorf("update member: %s", err)
	}
	return &updateResponse{
		Status: "Member updated",
	}, nil
}

func NewUpdateHandler(logger logr.Logger, db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	return handler.Create(logger, &updateHandler{db: db}, handler.WithOnlyForSuperuser(sessionsStore))
}
