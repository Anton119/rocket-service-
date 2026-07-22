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

// ShieldProperties — параметры щита.
type ShieldProperties struct {
	shieldType ShieldType
}

// ShieldType возвращает тип щита.
func (s *ShieldProperties) ShieldType() ShieldType { return s.shieldType }

// ConflictsWith проверяет конфликт щита с оружием (plasma + laser).
func (s *ShieldProperties) ConflictsWith(w *WeaponProperties) bool {
	if s == nil || w == nil {
		return false
	}
	return s.shieldType == ShieldTypePlasma && w.WeaponType() == WeaponTypeLaser
}

// NewShieldProperties создаёт свойства щита с валидацией типа.
func NewShieldProperties(shieldType ShieldType) (PartProperties, error) {
	switch shieldType {
	case ShieldTypeEnergy, ShieldTypePlasma:
		return NewPartProperties(nil, nil, &ShieldProperties{shieldType: shieldType}, nil), nil
	default:
		return PartProperties{}, fmt.Errorf("неизвестный тип щита %q: %w", shieldType, errs.ErrInvalidProperties)
	}
}
