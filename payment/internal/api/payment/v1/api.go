package v1

import (
	paymentv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/payment/v1"
)

// API — gRPC-адаптер PaymentService поверх доменного сервиса.
type API struct {
	paymentv1.UnimplementedPaymentServiceServer
	svc PaymentService
}

// NewAPI создаёт обработчик gRPC.
func NewAPI(svc PaymentService) *API {
	return &API{svc: svc}
}
