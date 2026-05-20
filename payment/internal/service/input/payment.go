package input

import (
	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/payment/internal/model"
)

// PayOrderInput — вход use case оплаты заказа.
type PayOrderInput struct {
	OrderUUID uuid.UUID
	Method    model.PaymentMethod
}

// PayOrderResult — результат оплаты.
type PayOrderResult struct {
	TransactionUUID uuid.UUID
}
