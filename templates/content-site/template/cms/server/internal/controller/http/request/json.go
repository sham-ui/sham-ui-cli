package request

import (
	"cms/pkg/tracing"
	"encoding/json"
	"fmt"
	"net/http"
)

func DecodeJSON[T any](r *http.Request) (*T, error) {
	const op = "DecodeJSON"

	_, span := tracing.StartSpan(r.Context(), scopeName, op)
	defer span.End()

	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode JSON: %w", err)
	}
	return &data, nil
}
