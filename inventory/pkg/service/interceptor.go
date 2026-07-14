package service

import (
	"google.golang.org/grpc"

	invinterceptor "github.com/Anton119/rocket-service-/inventory/internal/interceptor"
)

// UnaryErrorInterceptor переводит доменные ошибки inventory в gRPC status codes.
// Экспортируется через pkg для API-тестов и других модулей без доступа к internal.
func UnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return invinterceptor.UnaryErrorInterceptor()
}
