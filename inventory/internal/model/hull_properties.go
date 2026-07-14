package model

import (
	"fmt"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
)

const (
	minHullStrength = 30
	maxHullStrength = 200
)

// HullProperties — свойства корпуса.
type HullProperties struct {
	strength int
}

// NewHullProperties создаёт свойства корпуса с валидацией.
func NewHullProperties(strength int) (PartProperties, error) {
	if strength < minHullStrength || strength > maxHullStrength {
		return PartProperties{}, fmt.Errorf(
			"прочность корпуса %d вне диапазона [%d, %d]: %w",
			strength, minHullStrength, maxHullStrength, errs.ErrInvalidProperties,
		)
	}

	return PartProperties{hull: &HullProperties{strength: strength}}, nil
}

func (h *HullProperties) Strength() int { return h.strength }

// CanSupport проверяет, выдержит ли корпус нагрузку двигателя.
func (h *HullProperties) CanSupport(e *EngineProperties) bool {
	return h.strength >= e.requiredStrength
}
