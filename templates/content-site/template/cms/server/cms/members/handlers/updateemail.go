package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"net/http"
	"cms/core/handler"
	"cms/core/sessions"
	"cms/members"
	"strings"
)

type updateEmailHandler struct {
	db *sql.DB
}

type updateEmailRequestData struct {
	NewEmail1, NewEmail2 string
}

func (h *updateEmailHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	var data updateEmailRequestData
	err := json.NewDecoder(ctx.Request.Body).Decode(&data)
	if nil != err {
		return nil, fmt.Errorf("can't extract json data: %s", err)
	}
	data.NewEmail1 = strings.TrimSpace(data.NewEmail1)
	data.NewEmail2 = strings.TrimSpace(data.NewEmail2)
	return &data, nil
}

func (h *updateEmailHandler) Validate(ctx *handler.Context, data interface{}) (*handler.Validation, error) {
	validation := handler.NewValidation()
	requestData := data.(*updateEmailRequestData)
	if len(requestData.NewEmail1) <= 0 || len(requestData.NewEmail2) <= 0 {
		validation.AddError("Email must have more than 0 characters.")
	}
	if requestData.NewEmail1 != requestData.NewEmail2 {
		validation.AddError("Emails don't match.")
	}
	isUnique, err := members.IsUniqueEmail(h.db, requestData.NewEmail1)
	if nil != err {
		return nil, fmt.Errorf("is unique email: %s", err)
	}
	if !isUnique {
		validation.AddError("Email is already in use.")
	}
	return validation, nil
}

type updateEmailResponse struct {
	Status string `json:"Status"`
}

func (h *updateEmailHandler) Process(ctx *handler.Context, data interface{}) (interface{}, error) {
	requestData := data.(*updateEmailRequestData)
	session, _ := ctx.GetSession()
	err := members.UpdateMemberEmail(h.db, session.ID, requestData.NewEmail1)
	if nil != err {
		return nil, fmt.Errorf("can't update member email: %s", err)
	}
	rawSession := session.GetRawSession()
	rawSession.Values["email"] = requestData.NewEmail1
	if err = rawSession.Save(ctx.Request, ctx.Response); err != nil {
		return nil, fmt.Errorf("can't update session: %s", err)
	}
	return &updateEmailResponse{Status: "Email updated"}, nil
}

func NewUpdateEmailHandler(logger logr.Logger, db *sql.DB, sessionsStore *sessions.Store) http.HandlerFunc {
	return handler.Create(logger, &updateEmailHandler{db: db}, handler.WithOnlyForAuthenticated(sessionsStore))
}
