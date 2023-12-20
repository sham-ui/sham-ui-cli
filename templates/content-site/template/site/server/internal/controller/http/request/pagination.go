package request

import (
	"net/http"
	"site/pkg/tracing"
	"strconv"

	"site/internal/controller/http/response"
	"site/pkg/validation"
)

const scopeName = "http.request,pagination"

const (
	DefaultOffset = 0
	DefaultLimit  = 20
)

type PaginationParams struct {
	Offset int64
	Limit  int64
}

func ExtractPaginationParams(rw http.ResponseWriter, r *http.Request) (*PaginationParams, bool) {
	const op = "ExtractPaginationParams"

	_, span := tracing.StartSpan(r.Context(), scopeName, op)
	defer span.End()

	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		offset = DefaultOffset
	}
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		limit = DefaultLimit
	}
	requestValidation := validation.New()
	if offset < 0 {
		requestValidation.AddErrors("offset must be greater or equal than 0")
	}
	if limit <= 0 {
		requestValidation.AddErrors("limit must be greater than 0")
	}
	if requestValidation.HasErrors() {
		response.BadRequest(rw, r, requestValidation.Errors...)
		return nil, false
	}
	return &PaginationParams{
		Offset: offset,
		Limit:  limit,
	}, true
}
