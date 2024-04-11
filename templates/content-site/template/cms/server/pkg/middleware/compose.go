package middleware

import "net/http"

// Compose composes a chain of middleware functions into a single middleware function.
//
// It takes one or more middleware functions as parameters and returns a new middleware function.
// The returned middleware function takes an http.Handler as input and returns an http.Handler.
// When the returned middleware function is called, it executes each middleware function in passed order,
// passing the result of each middleware function as the next argument to the previous middleware function.
// The final middleware function in the chain is then returned.
func Compose(middleware ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		for i := len(middleware) - 1; i >= 0; i-- {
			next = middleware[i](next)
		}
		return next
	}
}
