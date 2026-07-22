package model

import (
	"fmt"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
)

// WeaponType — тип оружия.
type WeaponType string

const (
	WeaponTypeLaser   WeaponType = "laser"
	WeaponTypeMissile WeaponType = "missile"
)

// WeaponProperties — параметры оружия.
type WeaponProperties struct {
	weaponType WeaponType
}

// WeaponType возвращает тип оружия.
func (w *WeaponProperties) WeaponType() WeaponType { return w.weaponType }

// NewWeaponProperties создаёт свойства оружия с валидацией типа.
func NewWeaponProperties(weaponType WeaponType) (PartProperties, error) {
	switch weaponType {
	case WeaponTypeLaser, WeaponTypeMissile:
		return NewPartProperties(nil, nil, nil, &WeaponProperties{weaponType: weaponType}), nil
	default:
		return PartProperties{}, fmt.Errorf("неизвестный тип оружия %q: %w", weaponType, errs.ErrInvalidProperties)
	}
}
