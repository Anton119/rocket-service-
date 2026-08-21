package app

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	payapi "github.com/Anton119/rocket-service-/payment/internal/api/payment/v1"
	paymentsvc "github.com/Anton119/rocket-service-/payment/internal/service/payment"
	paymentv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/payment/v1"
)

func newGRPCServer(opts ...grpc.ServerOption) *grpc.Server {
	all := append([]grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}, opts...)

	return grpc.NewServer(all...)
}

type diContainer struct {
	paymentSvc *paymentsvc.Service
	api        paymentv1.PaymentServiceServer
}

func (d *diContainer) PaymentService() *paymentsvc.Service {
	if d.paymentSvc == nil {
		d.paymentSvc = paymentsvc.NewService()
	}
	return d.paymentSvc
}

func (d *diContainer) PaymentAPI() paymentv1.PaymentServiceServer {
	if d.api == nil {
		d.api = payapi.NewAPI(d.PaymentService())
	}
	return d.api
}
