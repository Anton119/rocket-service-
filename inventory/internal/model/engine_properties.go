package model

import (
	"fmt"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
)

// EngineClass — класс ионного двигателя.
type EngineClass string

const (
	EngineClassA EngineClass = "A"
	EngineClassB EngineClass = "B"
	EngineClassC EngineClass = "C"
)

// EngineProperties — свойства двигателя.
type EngineProperties struct {
	class            EngineClass
	requiredStrength int
}

// NewEngineProperties создаёт свойства двигателя с валидацией.
func NewEngineProperties(class EngineClass, requiredStrength int) (PartProperties, error) {
	switch class {
	case EngineClassA, EngineClassB, EngineClassC:
	default:
		return PartProperties{}, fmt.Errorf("неизвестный класс двигателя %q: %w", class, errs.ErrInvalidProperties)
	}

	if requiredStrength <= 0 {
		return PartProperties{}, fmt.Errorf("required_strength должен быть положительным: %w", errs.ErrInvalidProperties)
	}

	return PartProperties{
		engine: &EngineProperties{
			class:            class,
			requiredStrength: requiredStrength,
		},
	}, nil
}

func (e *EngineProperties) Class() EngineClass    { return e.class }
func (e *EngineProperties) RequiredStrength() int { return e.requiredStrength }
