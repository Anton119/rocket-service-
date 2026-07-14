package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	invinterceptor "github.com/Anton119/rocket-service-/inventory/internal/interceptor"
)

func withErrorInterceptor[T any](
	ctx context.Context,
	req any,
	fn func(context.Context, any) (any, error),
) (T, error) {
	var zero T

	interceptor := invinterceptor.UnaryErrorInterceptor()
	resp, err := interceptor(ctx, req, &grpc.UnaryServerInfo{}, fn)
	if err != nil {
		return zero, err
	}

	result, ok := resp.(T)
	if !ok {
		return zero, fmt.Errorf("неожиданный тип ответа: %T", resp)
	}

	return result, nil
}
