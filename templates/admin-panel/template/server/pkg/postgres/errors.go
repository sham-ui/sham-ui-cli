package postgres

import (
	"errors"

	"github.com/lib/pq"
)

func IsUniqueViolationError(err error, constraint string) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" && pgErr.Constraint == constraint
	}
	return false
}

func IsForeignKeyViolationError(err error, constraint string) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23503" && pgErr.Constraint == constraint
	}
	return false
}
