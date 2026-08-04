package converter

import (
	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/order/internal/model"
	"github.com/Anton119/rocket-service-/order/internal/service/input"
	orderv1 "github.com/Anton119/rocket-service-/shared/pkg/openapi/order/v1"
)

// PaymentMethodFromOpenAPI сопоставляет DTO способа оплаты доменному перечислению.
func PaymentMethodFromOpenAPI(m orderv1.PaymentMethod) (model.PaymentMethod, bool) {
	switch m {
	case orderv1.PaymentMethodCARD:
		return model.PaymentMethodCard, true
	case orderv1.PaymentMethodSBP:
		return model.PaymentMethodSBP, true
	case orderv1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard, true
	case orderv1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney, true
	default:
		return model.PaymentMethodInvalid, false
	}
}

// OrderToDTO собирает ответ GET /orders/{id}.
func OrderToDTO(order model.Order) *orderv1.OrderDto {
	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		shieldUUID = orderv1.NewOptNilUUID(*order.ShieldUUID)
	}

	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		weaponUUID = orderv1.NewOptNilUUID(*order.WeaponUUID)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	return &orderv1.OrderDto{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}

// CreateOrderInputFromRequest собирает вход use case из тела POST /orders.
func CreateOrderInputFromRequest(req *orderv1.CreateOrderRequest) input.CreateOrderInput {
	in := input.CreateOrderInput{
		HullUUID:   req.GetHullUUID(),
		EngineUUID: req.GetEngineUUID(),
	}
	if shield, ok := req.GetShieldUUID().Get(); ok {
		sh := shield
		in.ShieldUUID = &sh
	}
	if weapon, ok := req.WeaponUUID.Get(); ok {
		w := weapon
		in.WeaponUUID = &w
	}

	return in
}

// PayOrderInputFromRequest собирает вход use case из POST /orders/{order_uuid}/pay.
// Второй результат false, если способ оплаты в DTO неизвестен.
func PayOrderInputFromRequest(req *orderv1.PayOrderRequest, orderUUID uuid.UUID) (input.PayOrderInput, bool) {
	pm, ok := PaymentMethodFromOpenAPI(req.GetPaymentMethod())
	if !ok {
		return input.PayOrderInput{}, false
	}

	return input.PayOrderInput{
		OrderUUID:           orderUUID,
		Method:              pm,
		PaymentMethodStored: string(req.GetPaymentMethod()),
	}, true
}
