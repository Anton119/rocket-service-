package part

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
)

// ValidateCompatibility проверяет совместимость набора деталей по UUID.
func (s *Service) ValidateCompatibility(ctx context.Context, uuids []uuid.UUID) error {
	if err := validateUniqueUUIDs(uuids); err != nil {
		return err
	}

	parts, err := s.repo.ListParts(ctx, model.PartTypeUnspecified, uuids)
	if err != nil {
		return err
	}

	resolved, err := resolveShipParts(parts)
	if err != nil {
		return err
	}

	return s.compatibilityChecker.Check(resolved)
}

func validateUniqueUUIDs(uuids []uuid.UUID) error {
	seen := make(map[uuid.UUID]struct{}, len(uuids))
	for _, id := range uuids {
		if _, ok := seen[id]; ok {
			return errs.ErrPartTypeMismatch
		}
		seen[id] = struct{}{}
	}
	return nil
}

func resolveShipParts(parts []model.Part) (model.ResolvedShipSlots, error) {
	var (
		resolved model.ResolvedShipSlots
		hulls    int
		engines  int
		shields  int
		weapons  int
	)

	for _, p := range parts {
		switch p.PartType {
		case model.PartTypeHull:
			hulls++
			resolved.Hull = p
		case model.PartTypeEngine:
			engines++
			resolved.Engine = p
		case model.PartTypeShield:
			shields++
			part := p
			resolved.Shield = &part
		case model.PartTypeWeapon:
			weapons++
			part := p
			resolved.Weapon = &part
		default:
			return model.ResolvedShipSlots{}, errs.ErrPartTypeMismatch
		}
	}

	if hulls != 1 || engines != 1 || shields > 1 || weapons > 1 {
		return model.ResolvedShipSlots{}, errs.ErrPartTypeMismatch
	}

	return resolved, nil
}
