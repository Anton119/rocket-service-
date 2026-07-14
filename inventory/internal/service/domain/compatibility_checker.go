package domain

import (
	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
)

// CompatibilityChecker проверяет бизнес-правила совместимости деталей корабля.
type CompatibilityChecker struct{}

// NewCompatibilityChecker создаёт доменный сервис совместимости.
func NewCompatibilityChecker() *CompatibilityChecker {
	return &CompatibilityChecker{}
}

// Check проверяет бизнес-правила совместимости для набора деталей в слотах.
func (c *CompatibilityChecker) Check(slots model.ResolvedShipSlots) error {
	hullProps := slots.Hull.Properties().Hull()
	engineProps := slots.Engine.Properties().Engine()

	if hullProps == nil || engineProps == nil {
		return errs.ErrInvalidProperties
	}

	if !hullProps.CanSupport(engineProps) {
		return errs.ErrIncompatibleParts
	}

	if slots.Shield != nil && slots.Weapon != nil {
		shieldProps := slots.Shield.Properties().Shield()
		weaponProps := slots.Weapon.Properties().Weapon()

		if shieldProps == nil || weaponProps == nil {
			return errs.ErrInvalidProperties
		}

		if shieldProps.ConflictsWith(weaponProps) {
			return errs.ErrIncompatibleParts
		}
	}

	return nil
}
