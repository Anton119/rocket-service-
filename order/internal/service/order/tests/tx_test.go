package order_test

import "context"

type passthroughTx struct{}

func (passthroughTx) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
