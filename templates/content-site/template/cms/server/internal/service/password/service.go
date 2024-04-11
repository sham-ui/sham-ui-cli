package password

import (
	"cms/internal/model"
	"cms/pkg/tracing"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const scopeName = "service.password"

const passwordCost = 14

type service struct{}

func (s *service) Hash(ctx context.Context, raw string) (model.MemberHashedPassword, error) {
	const op = "Hash"

	_, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	bytes, err := bcrypt.GenerateFromPassword([]byte(raw), passwordCost)
	if err != nil {
		return bytes, fmt.Errorf("generate from password: %w", err)
	}
	return bytes, nil
}

func (s *service) Compare(ctx context.Context, hashed model.MemberHashedPassword, raw string) error {
	const op = "Compare"

	_, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	if err := bcrypt.CompareHashAndPassword(hashed, []byte(raw)); err != nil {
		return fmt.Errorf("compare hash and password: %w", err)
	}
	return nil
}

func New() *service {
	return &service{}
}
