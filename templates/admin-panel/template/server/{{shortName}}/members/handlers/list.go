package handlers

import (
	"database/sql"
	"fmt"
	"github.com/go-logr/logr"
	"net/http"
	"{{shortName}}/core/handler"
	"{{shortName}}/core/sessions"
	"strconv"
)

type listHandler struct {
	db *sql.DB
}

type listHandlerData struct {
	offset int
	limit  int
}

func (h *listHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	offset, err := strconv.Atoi(ctx.Request.URL.Query().Get("offset"))
	if nil != err {
		offset = 0
	}
	limit, err := strconv.Atoi(ctx.Request.URL.Query().Get("limit"))
	if nil != err {
		limit = 20
	}
	return &listHandlerData{
		offset: offset,
		limit:  limit,
	}, nil
}

func (h *listHandler) Validate(_ *handler.Context, data interface{}) (*handler.Validation, error) {
	validation := handler.NewValidation()
	params := data.(*listHandlerData)
	if params.limit <= 0 {
		validation.AddError("limit must be > 0")
	}
	if params.offset < 0 {
		validation.AddError("offset must be >= 0")
	}
	return validation, nil
}

func (h *listHandler) getMembersCount() (int, error) {
	var count int
	err := h.db.QueryRow("SELECT COUNT(*) FROM members").Scan(&count)
	if nil != err {
		return count, fmt.Errorf("select count: %s", err)
	}
	return count, nil
}

type listMemberData struct {
	ID, Name, Email string
	IsSuperuser     bool
}

func (h *listHandler) getMembers(offset, limit int) ([]*listMemberData, error) {
	var members []*listMemberData
	rows, err := h.db.Query("SELECT id, email, name, is_superuser FROM members ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
	defer rows.Close()
	if nil != err {
		return members, fmt.Errorf("query: %s", err)
	}
	for rows.Next() {
		data := &listMemberData{}
		err := rows.Scan(&data.ID, &data.Email, &data.Name, &data.IsSuperuser)
		if nil != err {
			return members, fmt.Errorf("scan row: %s", err)
		}
		members = append(members, data)
	}
	return members, nil
}

func (h *listHandler) Process(_ *handler.Context, data interface{}) (interface{}, error) {
	params := data.(*listHandlerData)
	count, err := h.getMembersCount()
	if nil != err {
		return nil, fmt.Errorf("members count: %s", err)
	}
	members, err := h.getMembers(params.offset, params.limit)
	if nil != err {
		return nil, fmt.Errorf("get members: %s", err)
	}
	return map[string]interface{}{
		"meta": map[string]int{
			"offset": params.offset,
			"limit":  params.limit,
			"total":  count,
		},
		"members": members,
	}, nil
}

func NewListHandler(logger logr.Logger, db *sql.DB, sessionStore *sessions.Store) http.HandlerFunc {
	return handler.Create(logger, &listHandler{db: db}, handler.WithOnlyForSuperuser(sessionStore))
}
