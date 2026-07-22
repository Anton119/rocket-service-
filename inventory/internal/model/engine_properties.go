package model

import (
	"fmt"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
)

// EngineClass — класс двигателя.
type EngineClass string

const (
	EngineClassA EngineClass = "A"
	EngineClassB EngineClass = "B"
	EngineClassC EngineClass = "C"
)

// EngineProperties — параметры двигателя.
type EngineProperties struct {
	class            EngineClass
	requiredStrength int
}

// Class возвращает класс двигателя.
func (e *EngineProperties) Class() EngineClass { return e.class }

// RequiredStrength возвращает минимальную прочность корпуса для двигателя.
func (e *EngineProperties) RequiredStrength() int { return e.requiredStrength }

// RequiredStrengthForClass возвращает требуемую прочность корпуса для класса двигателя.
func RequiredStrengthForClass(class EngineClass) int {
	switch class {
	case EngineClassA:
		return 100
	case EngineClassB:
		return 70
	case EngineClassC:
		return 30
	default:
		return 30
	}
}

// NewEngineProperties создаёт свойства двигателя с валидацией класса.
func NewEngineProperties(class EngineClass) (PartProperties, error) {
	switch class {
	case EngineClassA, EngineClassB, EngineClassC:
		return NewPartProperties(nil, &EngineProperties{
			class:            class,
			requiredStrength: RequiredStrengthForClass(class),
		}, nil, nil), nil
	default:
		return PartProperties{}, fmt.Errorf("неизвестный класс двигателя %q: %w", class, errs.ErrInvalidProperties)
	}
}
