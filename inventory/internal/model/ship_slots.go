package model

import "github.com/google/uuid"

// ShipSlots — именованные слоты корабля (для валидации совместимости по слотам).
type ShipSlots struct {
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID uuid.UUID
	WeaponUUID uuid.UUID
}
