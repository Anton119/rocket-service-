package input

import (
	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/order/internal/model"
)

// CreateOrderInput — вход use case создания заказа.
type CreateOrderInput struct {
	UserUUID   uuid.UUID
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID
}

// CreateOrderResult — результат создания заказа.
type CreateOrderResult struct {
	OrderUUID  uuid.UUID
	TotalPrice int64
}

// PayOrderInput — вход use case оплаты.
type PayOrderInput struct {
	OrderUUID           uuid.UUID
	Method              model.PaymentMethod
	PaymentMethodStored string // как в HTTP (enum string) для сохранения в модели
}
