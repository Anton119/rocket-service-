package converter

import (
	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/inventory/internal/model"
)

// ShipSlotsFromParts собирает слоты корабля по UUID деталей известных типов.
func ShipSlotsFromParts(hull, engine, shield, weapon uuid.UUID) model.ShipSlots {
	return model.ShipSlots{
		HullUUID:   hull,
		EngineUUID: engine,
		ShieldUUID: shield,
		WeaponUUID: weapon,
	}
}
