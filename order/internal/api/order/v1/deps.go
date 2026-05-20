package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/order/internal/model"
	"github.com/Anton119/rocket-service-/order/internal/service/input"
)

// OrderService описывает сервисный слой для HTTP API заказов.
type OrderService interface {
	CreateOrder(ctx context.Context, in input.CreateOrderInput) (*input.CreateOrderResult, error)
	GetOrder(ctx context.Context, id uuid.UUID) (model.Order, error)
	PayOrder(ctx context.Context, in input.PayOrderInput) (uuid.UUID, error)
	CancelOrder(ctx context.Context, id uuid.UUID) error
}
