package app

import (
	"context"

	payapi "github.com/Anton119/rocket-service-/payment/internal/api/payment/v1"
	paymentsvc "github.com/Anton119/rocket-service-/payment/internal/service/payment"
	paymentv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentSvc     payapi.PaymentService
	paymentHandler paymentv1.PaymentServiceServer
}

func (d *diContainer) PaymentService(_ context.Context) payapi.PaymentService {
	if d.paymentSvc == nil {
		d.paymentSvc = paymentsvc.NewService()
	}

	return d.paymentSvc
}

func (d *diContainer) PaymentV1API(ctx context.Context) paymentv1.PaymentServiceServer {
	if d.paymentHandler == nil {
		d.paymentHandler = payapi.NewAPI(d.PaymentService(ctx))
	}

	return d.paymentHandler
}
