package part

import (
	"context"

	"github.com/Anton119/rocket-service-/inventory/internal/model"
)

// ValidateCompatibility проверяет совместимость деталей в слотах корабля.
func (s *Service) ValidateCompatibility(ctx context.Context, slots model.ShipSlots) error {
	resolved, err := s.resolveShipSlots(ctx, slots)
	if err != nil {
		return err
	}

	return s.compatibilityChecker.Check(resolved)
}
