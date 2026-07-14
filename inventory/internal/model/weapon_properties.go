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

// WeaponProperties — свойства оружия.
type WeaponProperties struct {
	weaponType WeaponType
}

// NewWeaponProperties создаёт свойства оружия с валидацией.
func NewWeaponProperties(weaponType WeaponType) (PartProperties, error) {
	switch weaponType {
	case WeaponTypeLaser, WeaponTypeMissile:
	default:
		return PartProperties{}, fmt.Errorf("неизвестный тип оружия %q: %w", weaponType, errs.ErrInvalidProperties)
	}

	return PartProperties{weapon: &WeaponProperties{weaponType: weaponType}}, nil
}

func (w *WeaponProperties) WeaponType() WeaponType { return w.weaponType }
