package model

import (
	"fmt"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
)

const (
	minHullStrength = 30
	maxHullStrength = 200
)

// HullProperties — прочность корпуса.
type HullProperties struct {
	strength int
}

// Strength возвращает прочность корпуса.
func (h *HullProperties) Strength() int { return h.strength }

// CanSupport проверяет, выдержит ли корпус нагрузку двигателя.
func (h *HullProperties) CanSupport(e *EngineProperties) bool {
	if h == nil || e == nil {
		return false
	}
	return h.strength >= e.RequiredStrength()
}

// NewHullProperties создаёт свойства корпуса с валидацией диапазона.
func NewHullProperties(strength int) (PartProperties, error) {
	if strength < minHullStrength || strength > maxHullStrength {
		return PartProperties{}, fmt.Errorf(
			"прочность корпуса %d вне диапазона [%d, %d]: %w",
			strength, minHullStrength, maxHullStrength, errs.ErrInvalidProperties,
		)
	}
	return NewPartProperties(&HullProperties{strength: strength}, nil, nil, nil), nil
}
