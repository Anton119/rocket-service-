package interceptor

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Anton119/rocket-service-/platform/pkg/auth"
)

// SessionValidator проверяет сессию через IAM.
type SessionValidator interface {
	Whoami(ctx context.Context, sessionUUID string) (userUUID string, err error)
}

var publicMethods = map[string]bool{
	"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo":      true,
	"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": true,
}

// AuthInterceptor проверяет session-uuid из incoming metadata через IAM Whoami.
func AuthInterceptor(validator SessionValidator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "отсутствуют metadata")
		}

		values := md.Get(auth.SessionMetadataKey)
		if len(values) == 0 || values[0] == "" {
			return nil, status.Error(codes.Unauthenticated, "отсутствует session-uuid в metadata")
		}

		sessionUUID := values[0]
		userUUID, err := validator.Whoami(ctx, sessionUUID)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "сессия недействительна или истекла")
		}

		ctx = auth.WithSessionUUID(ctx, sessionUUID)
		ctx = auth.WithUserUUID(ctx, userUUID)

		slog.InfoContext(ctx, "gRPC: пользователь аутентифицирован", "method", info.FullMethod)

		return handler(ctx, req)
	}
}
