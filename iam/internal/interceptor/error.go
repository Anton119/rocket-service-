package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/Anton119/rocket-service-/iam/internal/errors"
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
	case errors.Is(err, errs.ErrInvalidLogin),
		errors.Is(err, errs.ErrWeakPassword),
		errors.Is(err, errs.ErrEmptyCredential),
		errors.Is(err, errs.ErrEmptySessionID),
		errors.Is(err, errs.ErrInvalidUUID):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, errs.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, errs.ErrUserAlreadyExists.Error())
	case errors.Is(err, errs.ErrInvalidCredentials),
		errors.Is(err, errs.ErrSessionNotFound):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, errs.ErrUserNotFound):
		return status.Error(codes.NotFound, errs.ErrUserNotFound.Error())
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
