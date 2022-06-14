package handler

import (
	"cms/core/sessions"
)

type handlerOptions struct {
	sessionStore          *sessions.Store
	onlyForAuthenticated  bool
	onlyForSuperuser      bool
	serializeResultToJSON bool
}

type Option interface {
	apply(options *handlerOptions)
}

// funcHandlerOption wraps a function that modifies handlerOptions into an
// implementation of the Option interface.
type funcHandlerOption struct {
	f func(options *handlerOptions)
}

func (fho *funcHandlerOption) apply(options *handlerOptions) {
	fho.f(options)
}

func newFuncHandlerOption(f func(options *handlerOptions)) *funcHandlerOption {
	return &funcHandlerOption{
		f: f,
	}
}

func defaultHandlerOptions() handlerOptions {
	return handlerOptions{
		sessionStore:          nil,
		onlyForAuthenticated:  false,
		onlyForSuperuser:      false,
		serializeResultToJSON: true,
	}
}

// WithOnlyForAuthenticated set checking for only auth user
// For non-auth user server return 403 HTTP status code
func WithOnlyForAuthenticated(store *sessions.Store) Option {
	return newFuncHandlerOption(func(options *handlerOptions) {
		options.onlyForAuthenticated = true
		options.sessionStore = store
	})
}

// WithOnlyForAuthenticated set checking for only superuser
// For non superuser server return 403 HTTP status code
func WithOnlyForSuperuser(store *sessions.Store) Option {
	return newFuncHandlerOption(func(options *handlerOptions) {
		options.onlyForAuthenticated = true
		options.onlyForSuperuser = true
		options.sessionStore = store
	})
}

// WithSessionStore add session store dependency to handler
func WithSessionStore(store *sessions.Store) Option {
	return newFuncHandlerOption(func(options *handlerOptions) {
		options.sessionStore = store
	})
}

// WithoutSerializeResultToJSON disable serialization Process response to JSON
func WithoutSerializeResultToJSON() Option {
	return newFuncHandlerOption(func(options *handlerOptions) {
		options.serializeResultToJSON = false
	})
}
