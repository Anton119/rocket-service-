package order

import (
	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
	"github.com/Anton119/rocket-service-/order/internal/service/input"
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

func partUUIDsFromInput(in input.CreateOrderInput) []string {
	uuids := []string{in.HullUUID.String(), in.EngineUUID.String()}
	if in.ShieldUUID != nil {
		uuids = append(uuids, in.ShieldUUID.String())
	}
	if in.WeaponUUID != nil {
		uuids = append(uuids, in.WeaponUUID.String())
	}

	return uuids
}

func indexPartsByUUID(parts []model.Part) map[string]model.Part {
	byUUID := make(map[string]model.Part, len(parts))
	for i := range parts {
		p := parts[i]
		byUUID[p.UUID] = p
	}

	return byUUID
}

func validatePartsStock(byUUID map[string]model.Part, uuids []string) error {
	for _, id := range uuids {
		p, ok := byUUID[id]
		if !ok {
			return errs.ErrPartNotFound
		}
		if p.StockQuantity <= 0 {
			return errs.ErrPartOutOfStock
		}
	}

	return nil
}
