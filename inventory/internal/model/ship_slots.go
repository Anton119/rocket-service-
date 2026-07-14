package model

import "github.com/google/uuid"

// ShipSlots — именованные слоты корабля для проверки совместимости.
type ShipSlots struct {
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID
}

// ResolvedShipSlots — детали, привязанные к слотам после проверки типов.
type ResolvedShipSlots struct {
	Hull   Part
	Engine Part
	Shield *Part
	Weapon *Part
}
