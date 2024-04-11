package session

import (
	"github.com/gorilla/sessions"
)

type valueNotFoundError struct {
	key string
}

func (e valueNotFoundError) Error() string {
	return "value not found: " + e.key
}

func (e valueNotFoundError) Is(target error) bool {
	//nolint:errorlint
	_, ok := target.(valueNotFoundError)
	return ok
}

type valueTypeError struct {
	key string
}

func (e valueTypeError) Error() string {
	return "value type error: " + e.key
}

func (e valueTypeError) Is(target error) bool {
	//nolint:errorlint
	_, ok := target.(valueTypeError)
	return ok
}

func newValueNotFoundError(key string) error {
	return valueNotFoundError{key: key}
}

func newValueTypeError(key string) error {
	return valueTypeError{key: key}
}

type sessionValue interface {
	~int | ~string | ~bool
}

func getValue[T sessionValue](s *sessions.Session, key string) (T, error) { //nolint:ireturn
	var zeroValue T
	value, ok := s.Values[key]
	if !ok {
		return zeroValue, newValueNotFoundError(key)
	}
	typedValue, ok := value.(T)
	if !ok {
		return zeroValue, newValueTypeError(key)
	}
	return typedValue, nil
}
