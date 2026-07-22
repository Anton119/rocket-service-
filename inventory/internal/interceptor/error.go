package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
)

// ToGRPCError переводит доменные ошибки в gRPC status.
func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := status.FromError(err); ok {
		return err
	}

	switch {
	case errors.Is(err, errs.ErrPartNotFound):
		return status.Error(codes.NotFound, errs.ErrPartNotFound.Error())
	case errs.IsInvalidUUID(err):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, errs.ErrPartTypeMismatch):
		return status.Error(codes.InvalidArgument, errs.ErrPartTypeMismatch.Error())
	case errors.Is(err, errs.ErrIncompatibleParts):
		return status.Error(codes.FailedPrecondition, errs.ErrIncompatibleParts.Error())
	case errors.Is(err, errs.ErrOutOfStock):
		return status.Error(codes.ResourceExhausted, errs.ErrOutOfStock.Error())
	case errors.Is(err, errs.ErrNothingToRelease):
		return status.Error(codes.FailedPrecondition, errs.ErrNothingToRelease.Error())
	case errors.Is(err, errs.ErrNothingToCommit):
		return status.Error(codes.FailedPrecondition, errs.ErrNothingToCommit.Error())
	case errors.Is(err, errs.ErrInvalidProperties):
		return status.Error(codes.Internal, errs.ErrInvalidProperties.Error())
	default:
		return status.Errorf(codes.Internal, "внутренняя ошибка: %v", err)
	}
}

// ErrorInterceptor маппит доменные ошибки обработчиков в gRPC-коды.
func ErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return nil, ToGRPCError(err)
		}
		return resp, nil
	}
}
