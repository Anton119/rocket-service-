package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/Anton119/rocket-service-/platform/pkg/auth"
)

// SessionForwarder прокидывает session-uuid из context в outgoing gRPC metadata.
func SessionForwarder() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if sessionUUID, ok := auth.SessionUUIDFromContext(ctx); ok && sessionUUID != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, auth.SessionMetadataKey, sessionUUID)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
