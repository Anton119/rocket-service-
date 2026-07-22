package model

// ResolvedShipSlots — набор деталей корабля после проверки типов.
type ResolvedShipSlots struct {
	Hull   Part
	Engine Part
	Shield *Part
	Weapon *Part
}
