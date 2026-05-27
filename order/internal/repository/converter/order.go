package converter

import (
	"github.com/Anton119/rocket-service-/order/internal/model"
	"github.com/Anton119/rocket-service-/order/internal/repository/record"
)

func OrderToRecord(o model.Order) record.Order {
	return record.Order{
		OrderUUID:       o.OrderUUID,
		HullUUID:        o.HullUUID,
		EngineUUID:      o.EngineUUID,
		ShieldUUID:      o.ShieldUUID,
		WeaponUUID:      o.WeaponUUID,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   o.PaymentMethod,
		Status:          string(o.Status),
		CreatedAt:       o.CreatedAt,
	}
}

func OrderToModel(r record.Order) model.Order {
	return model.Order{
		OrderUUID:       r.OrderUUID,
		HullUUID:        r.HullUUID,
		EngineUUID:      r.EngineUUID,
		ShieldUUID:      r.ShieldUUID,
		WeaponUUID:      r.WeaponUUID,
		TotalPrice:      r.TotalPrice,
		TransactionUUID: r.TransactionUUID,
		PaymentMethod:   r.PaymentMethod,
		Status:          model.OrderStatus(r.Status),
		CreatedAt:       r.CreatedAt,
	}
}
