package ssr

import "errors"

var (
	errSSRProcessStopped = errors.New("ssr process stopped")
)

type (
	ServerSideRenderRequestError struct {
		Err error
	}
	ServerSideRenderRespondError struct {
		Message string
	}
)

func (e ServerSideRenderRequestError) Error() string {
	return "ssr request: " + e.Err.Error()
}

func (e ServerSideRenderRequestError) Unwrap() error {
	return e.Err
}

func (e ServerSideRenderRequestError) Is(target error) bool {
	//nolint:errorlint
	_, ok := target.(ServerSideRenderRequestError)
	return ok
}

func (e ServerSideRenderRespondError) Error() string {
	return "ssr respond with error: " + e.Message
}

func (e ServerSideRenderRespondError) Is(target error) bool {
	//nolint:errorlint
	_, ok := target.(ServerSideRenderRespondError)
	return ok
}

func NewServerSideRenderRequestError(err error) ServerSideRenderRequestError {
	return ServerSideRenderRequestError{
		Err: err,
	}
}

func NewServerSideRenderError(message string) ServerSideRenderRespondError {
	return ServerSideRenderRespondError{
		Message: message,
	}
}
