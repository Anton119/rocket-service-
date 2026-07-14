package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
	"github.com/Anton119/rocket-service-/inventory/internal/service/input"
)

func (s *Service) resolveShipSlots(ctx context.Context, slots model.ShipSlots) (model.ResolvedShipSlots, error) {
	ids := make([]uuid.UUID, 0, 4)
	ids = append(ids, slots.HullUUID, slots.EngineUUID)

	if slots.ShieldUUID != nil {
		ids = append(ids, *slots.ShieldUUID)
	}
	if slots.WeaponUUID != nil {
		ids = append(ids, *slots.WeaponUUID)
	}

	if hasDuplicateUUIDs(ids) {
		return model.ResolvedShipSlots{}, errs.ErrPartTypeMismatch
	}

	parts, err := s.repo.List(ctx, input.PartFilter{UUIDs: ids})
	if err != nil {
		return model.ResolvedShipSlots{}, fmt.Errorf("получить детали: %w", err)
	}

	byID := make(map[uuid.UUID]model.Part, len(parts))
	for _, p := range parts {
		byID[p.UUID()] = p
	}

	hull, ok := byID[slots.HullUUID]
	if !ok || hull.PartType() != model.PartTypeHull {
		return model.ResolvedShipSlots{}, errs.ErrPartTypeMismatch
	}

	engine, ok := byID[slots.EngineUUID]
	if !ok || engine.PartType() != model.PartTypeEngine {
		return model.ResolvedShipSlots{}, errs.ErrPartTypeMismatch
	}

	resolved := model.ResolvedShipSlots{
		Hull:   hull,
		Engine: engine,
	}

	if slots.ShieldUUID != nil {
		shield, shieldOK := byID[*slots.ShieldUUID]
		if !shieldOK || shield.PartType() != model.PartTypeShield {
			return model.ResolvedShipSlots{}, errs.ErrPartTypeMismatch
		}

		resolved.Shield = &shield
	}

	if slots.WeaponUUID != nil {
		weapon, weaponOK := byID[*slots.WeaponUUID]
		if !weaponOK || weapon.PartType() != model.PartTypeWeapon {
			return model.ResolvedShipSlots{}, errs.ErrPartTypeMismatch
		}

		resolved.Weapon = &weapon
	}

	return resolved, nil
}

func hasDuplicateUUIDs(ids []uuid.UUID) bool {
	seen := make(map[uuid.UUID]struct{}, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			return true
		}
		seen[id] = struct{}{}
	}

	return false
}
