package converter

import (
	"github.com/Anton119/rocket-service-/order/internal/model"
	"github.com/Anton119/rocket-service-/order/internal/repository/record"
)

func OrderToRecord(o model.Order) record.Order {
	return record.Order{
		UUID:            o.OrderUUID,
		UserUUID:        o.UserUUID,
		Status:          string(o.Status),
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   o.PaymentMethod,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}
}

func OrderItemsToRecord(o model.Order) []record.OrderItem {
	items := make([]record.OrderItem, len(o.Items))
	for i, item := range o.Items {
		items[i] = record.OrderItem{
			OrderUUID: o.OrderUUID,
			PartUUID:  item.PartUUID,
			PartType:  item.PartType,
			Price:     item.Price,
		}
	}
	return items
}

func OrderToModel(rec record.Order, items []record.OrderItem) model.Order {
	domainItems := make([]model.OrderItem, len(items))
	var total int64

	o := model.Order{
		OrderUUID:       rec.UUID,
		UserUUID:        rec.UserUUID,
		Status:          model.OrderStatus(rec.Status),
		TransactionUUID: rec.TransactionUUID,
		PaymentMethod:   rec.PaymentMethod,
		CreatedAt:       rec.CreatedAt,
		UpdatedAt:       rec.UpdatedAt,
	}

	for i, item := range items {
		domainItems[i] = model.OrderItem{
			PartUUID: item.PartUUID,
			PartType: item.PartType,
			Price:    item.Price,
		}
		total += item.Price

		switch item.PartType {
		case "HULL":
			o.HullUUID = item.PartUUID
		case "ENGINE":
			o.EngineUUID = item.PartUUID
		case "SHIELD":
			partUUID := item.PartUUID
			o.ShieldUUID = &partUUID
		case "WEAPON":
			partUUID := item.PartUUID
			o.WeaponUUID = &partUUID
		}
	}

	o.Items = domainItems
	o.TotalPrice = total

	return o
}
