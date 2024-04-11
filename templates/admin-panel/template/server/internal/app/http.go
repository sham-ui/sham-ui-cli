package app

import "{{ shortName }}/internal/controller/http"

type httpServerDependencies struct {
	passwordService http.PasswordService
	sessionService  http.SessionService
	memberService   http.MemberService
}

func (r httpServerDependencies) PasswordService() http.PasswordService { //nolint:ireturn
	return r.passwordService
}

func (r httpServerDependencies) SessionService() http.SessionService { //nolint:ireturn
	return r.sessionService
}

func (r httpServerDependencies) MemberService() http.MemberService { //nolint:ireturn
	return r.memberService
}
