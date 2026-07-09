package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/order/internal/model"
)

// InventoryClient — получение информации о деталях из каталога.
type InventoryClient interface {
	ListParts(ctx context.Context, uuids []string) ([]model.Part, error)
}

// PaymentClient — проведение оплаты во внешнем сервисе.
type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID string, method model.PaymentMethod) (uuid.UUID, error)
}

// OrderRepository — персистентность заказов.
type OrderRepository interface {
	Create(ctx context.Context, order model.Order) error
	Get(ctx context.Context, id uuid.UUID) (model.Order, error)
	Save(ctx context.Context, order model.Order) error
}
