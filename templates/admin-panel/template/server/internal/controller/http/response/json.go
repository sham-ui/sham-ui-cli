package response

import (
	"encoding/json"
	"net/http"
	"{{ shortName }}/pkg/logger"
	"{{ shortName }}/pkg/tracing"
)

const scopeName = "http.response"

type responseError struct {
	Status   string   `json:"Status"`
	Messages []string `json:"Messages,omitempty"`
}

func InternalServerError(rw http.ResponseWriter, r *http.Request) {
	WithError(rw, r, http.StatusInternalServerError, "internal server error")
}

func NotFound(rw http.ResponseWriter, r *http.Request, messages ...string) {
	WithError(rw, r, http.StatusNotFound, messages...)
}

func BadRequest(rw http.ResponseWriter, r *http.Request, messages ...string) {
	WithError(rw, r, http.StatusBadRequest, messages...)
}

func JSON(rw http.ResponseWriter, r *http.Request, statusCode int, msg any) {
	const op = "JSON"

	_, span := tracing.StartSpan(r.Context(), scopeName, op)
	defer span.End()

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	if err := json.NewEncoder(rw).Encode(msg); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		logger.Load(r.Context()).Error(err, "encode json message fail")
	}
}

func WithError(rw http.ResponseWriter, r *http.Request, statusCode int, messages ...string) {
	msg := &responseError{
		Status:   http.StatusText(statusCode),
		Messages: messages,
	}
	JSON(rw, r, statusCode, msg)
}
