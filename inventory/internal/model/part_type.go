package model

// PartType — тип детали в доменной модели (числовые значения совпадают с inventory.v1.PartType).
type PartType int32

const (
	PartTypeUnspecified PartType = 0
	PartTypeHull        PartType = 1
	PartTypeEngine      PartType = 2
	PartTypeShield      PartType = 3
	PartTypeWeapon      PartType = 4
)
