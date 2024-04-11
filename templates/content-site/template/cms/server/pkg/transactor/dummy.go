package transactor

import "context"

type dummyTransactor struct{}

func (dt dummyTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func NewDummyTransactor() dummyTransactor {
	return dummyTransactor{}
}
