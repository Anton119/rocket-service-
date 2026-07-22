package service

import (
	"google.golang.org/grpc"

	invinterceptor "github.com/Anton119/rocket-service-/inventory/internal/interceptor"
)

// UnaryErrorInterceptor маппит доменные ошибки InventoryService в gRPC-коды.
// Нужен снаружи модуля (API-тесты order), поэтому экспортируется через pkg.
func UnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return invinterceptor.ErrorInterceptor()
}
