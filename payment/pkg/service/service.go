package service

import (
	payapi "github.com/Anton119/rocket-service-/payment/internal/api/payment/v1"
	paymentsvc "github.com/Anton119/rocket-service-/payment/internal/service/payment"
	paymentv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/payment/v1"
)

// NewPaymentServer собирает gRPC PaymentService.
// Пакет pkg нужен модулям вроде order/tests, которые не могут импортировать payment/internal.
func NewPaymentServer() paymentv1.PaymentServiceServer {
	svc := paymentsvc.NewService()

	return payapi.NewAPI(svc)
}
