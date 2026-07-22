package domain

import (
	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
)

const defaultHullStrength = 100

// CompatibilityChecker проверяет совместимость деталей космического корабля.
type CompatibilityChecker struct{}

// NewCompatibilityChecker создаёт доменный сервис проверки совместимости.
func NewCompatibilityChecker() *CompatibilityChecker {
	return &CompatibilityChecker{}
}

// Check проверяет бизнес-правила совместимости для набора деталей.
func (c *CompatibilityChecker) Check(slots model.ResolvedShipSlots) error {
	hull := hullProperties(slots.Hull)
	engine := engineProperties(slots.Engine)
	if hull == nil || engine == nil {
		return errs.ErrIncompatibleParts
	}

	if !hull.CanSupport(engine) {
		return errs.ErrIncompatibleParts
	}

	if slots.Shield != nil && slots.Weapon != nil {
		shield := shieldProperties(*slots.Shield)
		weapon := weaponProperties(*slots.Weapon)
		if shield != nil && weapon != nil && shield.ConflictsWith(weapon) {
			return errs.ErrIncompatibleParts
		}
	}

	return nil
}

func hullProperties(p model.Part) *model.HullProperties {
	if h := p.Properties.Hull(); h != nil {
		return h
	}
	if p.PartType != model.PartTypeHull {
		return nil
	}
	props, err := model.NewHullProperties(defaultHullStrength)
	if err != nil {
		return nil
	}
	return props.Hull()
}

func engineProperties(p model.Part) *model.EngineProperties {
	if e := p.Properties.Engine(); e != nil {
		return e
	}
	if p.PartType != model.PartTypeEngine {
		return nil
	}
	props, err := model.NewEngineProperties(model.EngineClassB)
	if err != nil {
		return nil
	}
	return props.Engine()
}

func shieldProperties(p model.Part) *model.ShieldProperties {
	if s := p.Properties.Shield(); s != nil {
		return s
	}
	if p.PartType != model.PartTypeShield {
		return nil
	}
	props, err := model.NewShieldProperties(model.ShieldTypeEnergy)
	if err != nil {
		return nil
	}
	return props.Shield()
}

func weaponProperties(p model.Part) *model.WeaponProperties {
	if w := p.Properties.Weapon(); w != nil {
		return w
	}
	if p.PartType != model.PartTypeWeapon {
		return nil
	}
	props, err := model.NewWeaponProperties(model.WeaponTypeMissile)
	if err != nil {
		return nil
	}
	return props.Weapon()
}
