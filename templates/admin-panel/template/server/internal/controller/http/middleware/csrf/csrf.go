package csrf

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func New(authKey []byte, requestHeader string, cookieName string) func(next http.Handler) http.Handler {
	return csrf.Protect(
		authKey,
		csrf.RequestHeader(requestHeader),
		csrf.CookieName(cookieName),
		csrf.Secure(false),
	)
}
