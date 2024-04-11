package http

import (
	membersCreate "{{ shortName }}/internal/controller/http/handler/api/admin/members/create"
	membersList "{{ shortName }}/internal/controller/http/handler/api/admin/members/list"
	membersDelete "{{ shortName }}/internal/controller/http/handler/api/admin/members/remove"
	membersResetPassword "{{ shortName }}/internal/controller/http/handler/api/admin/members/reset_password"
	membersUpdate "{{ shortName }}/internal/controller/http/handler/api/admin/members/update"
	"{{ shortName }}/internal/controller/http/handler/api/members/email"
	"{{ shortName }}/internal/controller/http/handler/api/members/name"
	"{{ shortName }}/internal/controller/http/handler/api/members/password"
	{{#if signupEnabled}}
	"{{ shortName }}/internal/controller/http/handler/api/members/signup"
	{{/if}}
	"{{ shortName }}/internal/controller/http/handler/api/session/login"
	"{{ shortName }}/internal/controller/http/handler/api/session/logout"
	sessionMW "{{ shortName }}/internal/controller/http/middleware/session"
)

//nolint:lll
//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name HandlerDependencyProvider --inpackage --testonly --with-expecter
type HandlerDependencyProvider interface {
	PasswordService() PasswordService
	SessionService() SessionService
	MemberService() MemberService
}

type PasswordService interface {
	login.PasswordService
	{{#if signupEnabled}}
	signup.PasswordService
	{{/if}}
	password.PasswordService
	membersCreate.PasswordService
	membersResetPassword.PasswordService
}

type SessionService interface {
	sessionMW.Service
	login.SessionService
	logout.Service
	email.SessionService
	name.SessionService
}

type MemberService interface {
	login.MemberService
	{{#if signupEnabled}}
	signup.MemberService
	{{/if}}
	email.MemberService
	name.MemberService
	password.MemberService
	membersList.MemberService
	membersCreate.MemberService
	membersUpdate.MemberService
	membersResetPassword.MemberService
	membersDelete.MemberService
}
