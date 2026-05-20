package v1

import (
	"context"

	"github.com/Anton119/rocket-service-/payment/internal/service/input"
)

// PaymentService описывает сервисный слой для gRPC API оплаты.
type PaymentService interface {
	PayOrder(ctx context.Context, in input.PayOrderInput) (*input.PayOrderResult, error)
}
