package order

import (
	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/order/internal/model"
)

type orderSlot struct {
	partUUID uuid.UUID
	partType string
	price    int64
}

func mergeOrderItems(slots []orderSlot) []model.OrderItem {
	items := make([]model.OrderItem, 0, len(slots))
	indexByPart := make(map[uuid.UUID]int, len(slots))

	for _, slot := range slots {
		if idx, ok := indexByPart[slot.partUUID]; ok {
			items[idx].Price += slot.price
			continue
		}

		indexByPart[slot.partUUID] = len(items)
		items = append(items, model.OrderItem{
			PartUUID: slot.partUUID,
			PartType: slot.partType,
			Price:    slot.price,
		})
	}

	return items
}

func appendMergedItem(items []model.OrderItem, partUUID uuid.UUID, partType string, price int64) []model.OrderItem {
	for i := range items {
		if items[i].PartUUID == partUUID {
			items[i].Price += price
			return items
		}
	}

	return append(items, model.OrderItem{
		PartUUID: partUUID,
		PartType: partType,
		Price:    price,
	})
}
