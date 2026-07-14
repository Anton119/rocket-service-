package model

import (
	"fmt"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
)

// ShieldType — тип щита.
type ShieldType string

const (
	ShieldTypeEnergy ShieldType = "energy"
	ShieldTypePlasma ShieldType = "plasma"
)

// ShieldProperties — свойства щита.
type ShieldProperties struct {
	shieldType ShieldType
}

// NewShieldProperties создаёт свойства щита с валидацией.
func NewShieldProperties(shieldType ShieldType) (PartProperties, error) {
	switch shieldType {
	case ShieldTypeEnergy, ShieldTypePlasma:
	default:
		return PartProperties{}, fmt.Errorf("неизвестный тип щита %q: %w", shieldType, errs.ErrInvalidProperties)
	}

	return PartProperties{shield: &ShieldProperties{shieldType: shieldType}}, nil
}

func (s *ShieldProperties) ShieldType() ShieldType { return s.shieldType }

// ConflictsWith проверяет конфликт щита с оружием (plasma + laser).
func (s *ShieldProperties) ConflictsWith(w *WeaponProperties) bool {
	return s.shieldType == ShieldTypePlasma && w.WeaponType() == WeaponTypeLaser
}
