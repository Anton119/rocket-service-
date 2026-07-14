package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
)

// UnaryErrorInterceptor переводит доменные ошибки в gRPC status codes.
func UnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		return resp, mapError(err)
	}
}

func mapError(err error) error {
	switch {
	case errors.Is(err, errs.ErrPartNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, errs.ErrInvalidUUID), errors.Is(err, errs.ErrEmptyUUID), errors.Is(err, errs.ErrPartTypeMismatch):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, errs.ErrIncompatibleParts), errors.Is(err, errs.ErrNothingToRelease):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, errs.ErrOutOfStock):
		return status.Error(codes.ResourceExhausted, err.Error())
	case errors.Is(err, errs.ErrInvalidProperties):
		return status.Error(codes.Internal, err.Error())
	default:
		return err
	}
}
