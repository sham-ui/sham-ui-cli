package handlers

import (
	"cms/config"
	"cms/core/handler"
	"cms/core/sessions"
	"cms/members"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/gorilla/csrf"
	gorillaSessions "github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type loginHandler struct {
	sessionsStore *sessions.Store
	db            *sql.DB
}

// credentials are the user provided email and password
type credentials struct {
	Email, Password string
}

type loginRequestData struct {
	form    *credentials
	member  *members.MemberData
	session *gorillaSessions.Session
}

// decodeCredentials decodes the JSON data into a struct containing the email and password.DecodeCredentials
func (h *loginHandler) decodeCredentials(r *http.Request) (*credentials, error) {
	var c credentials
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		return nil, fmt.Errorf("decode: %s", err)
	}
	return &c, nil
}

func (h *loginHandler) getMemberData(email string) (*members.MemberData, error) {
	data := &members.MemberData{}
	err := h.db.QueryRow("SELECT id, email, name, password, is_superuser FROM members WHERE email =$1", email).Scan(&data.ID, &data.Email, &data.Name, &data.Password, &data.IsSuperuser)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("select member: %s", err)
	}
	return data, nil
}

func (h *loginHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	session, err := h.sessionsStore.GetRawSession(ctx.Request)
	if nil != err {
		return nil, fmt.Errorf("can't get session: %s", err)
	}
	cred, err := h.decodeCredentials(ctx.Request)
	if nil != err {
		return nil, fmt.Errorf("decode credentials: %s", err)
	}
	member, err := h.getMemberData(cred.Email)
	if nil != err {
		return nil, fmt.Errorf("get member: %s", err)
	}
	if nil == member {
		return nil, fmt.Errorf("member not found")
	}
	return &loginRequestData{
		form:    cred,
		member:  member,
		session: session,
	}, nil
}

func (h *loginHandler) Validate(_ *handler.Context, data interface{}) (*handler.Validation, error) {
	validation := handler.NewValidation()
	requestData := data.(*loginRequestData)
	if nil != bcrypt.CompareHashAndPassword([]byte(requestData.member.Password), []byte(requestData.form.Password)) {
		validation.AddError("Incorrect username or password")
	}
	if auth, ok := requestData.session.Values["authenticated"].(bool); ok && auth {
		validation.AddError("already logged")
	}
	return validation, nil
}

type loginResponse struct {
	Status      string
	Name        string
	Email       string
	IsSuperuser bool
}

func (h *loginHandler) Process(ctx *handler.Context, data interface{}) (interface{}, error) {
	requestData := data.(*loginRequestData)

	// Limit the sessions to 1 24-hour day
	requestData.session.Options.MaxAge = 86400 * 1
	requestData.session.Options.Domain = config.Server.AllowedDomains[0]
	requestData.session.Options.HttpOnly = true
	// Set cookie values and save
	requestData.session.Values["authenticated"] = true
	requestData.session.Values["name"] = requestData.member.Name
	requestData.session.Values["email"] = requestData.member.Email
	requestData.session.Values["id"] = requestData.member.ID
	requestData.session.Values["is_superuser"] = requestData.member.IsSuperuser
	if err := requestData.session.Save(ctx.Request, ctx.Response); err != nil {
		return nil, fmt.Errorf("save session: %s", err)
	}
	ctx.Response.Header().Set("X-CSRF-Token", csrf.Token(ctx.Request))
	return &loginResponse{
		Status:      "OK",
		Name:        requestData.member.Name,
		Email:       requestData.member.Email,
		IsSuperuser: requestData.member.IsSuperuser,
	}, nil
}

func NewLoginHandler(logger logr.Logger, sessionsStore *sessions.Store, db *sql.DB) http.HandlerFunc {
	return handler.Create(logger, &loginHandler{
		sessionsStore: sessionsStore,
		db:            db,
	})
}
