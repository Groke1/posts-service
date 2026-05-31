package db

import "context"

type NoopTransactor struct{}

func NewNoopTransactor() *NoopTransactor {
	return &NoopTransactor{}
}

func (t *NoopTransactor) WithTx(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return fn(ctx)
}
