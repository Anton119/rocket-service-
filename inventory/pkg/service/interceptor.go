package service

import (
	"context"

	"google.golang.org/grpc"

	invinterceptor "github.com/Anton119/rocket-service-/inventory/internal/interceptor"
)

// UnaryErrorInterceptor маппит доменные ошибки InventoryService в gRPC-коды.
// Нужен снаружи модуля (API-тесты order), поэтому экспортируется через pkg.
func UnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return invinterceptor.ErrorInterceptor()
}

// SessionValidator проверяет сессию через IAM Whoami.
type SessionValidator interface {
	Whoami(ctx context.Context, sessionUUID string) (userUUID string, err error)
}

// UnaryAuthInterceptor проверяет session-uuid из incoming metadata.
func UnaryAuthInterceptor(validator SessionValidator) grpc.UnaryServerInterceptor {
	return invinterceptor.AuthInterceptor(validator)
}
