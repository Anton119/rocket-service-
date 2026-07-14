package converter

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
	"github.com/Anton119/rocket-service-/inventory/internal/service/input"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

// ListPartsInputFromRequest собирает вход use case из rpc ListParts.
func ListPartsInputFromRequest(req *inventoryv1.ListPartsRequest) (input.ListPartsInput, error) {
	in := input.ListPartsInput{
		PartType: PartTypeFromProto(req.GetPartType()),
	}
	if len(req.GetUuids()) == 0 {
		return in, nil
	}

	ids := make([]uuid.UUID, 0, len(req.GetUuids()))
	for _, idStr := range req.GetUuids() {
		id, err := ParsePartUUID(idStr)
		if err != nil {
			return input.ListPartsInput{}, err
		}
		ids = append(ids, id)
	}
	in.IDs = ids

	return in, nil
}

// ShipSlotsFromRequest собирает слоты корабля из rpc ValidateCompatibility.
func ShipSlotsFromRequest(req *inventoryv1.ValidateCompatibilityRequest) (model.ShipSlots, error) {
	hullUUID, err := ParsePartUUID(req.GetHullUuid())
	if err != nil {
		return model.ShipSlots{}, err
	}

	engineUUID, err := ParsePartUUID(req.GetEngineUuid())
	if err != nil {
		return model.ShipSlots{}, err
	}

	slots := model.ShipSlots{
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
	}

	if req.GetShieldUuid() != "" {
		shieldUUID, parseErr := ParsePartUUID(req.GetShieldUuid())
		if parseErr != nil {
			return model.ShipSlots{}, parseErr
		}
		slots.ShieldUUID = &shieldUUID
	}

	if req.GetWeaponUuid() != "" {
		weaponUUID, parseErr := ParsePartUUID(req.GetWeaponUuid())
		if parseErr != nil {
			return model.ShipSlots{}, parseErr
		}
		slots.WeaponUUID = &weaponUUID
	}

	return slots, nil
}

// UUIDsFromStrings разбирает список UUID из rpc Reserve/Release.
func UUIDsFromStrings(ids []string) ([]uuid.UUID, error) {
	out := make([]uuid.UUID, 0, len(ids))
	for _, idStr := range ids {
		id, err := ParsePartUUID(idStr)
		if err != nil {
			return nil, err
		}
		out = append(out, id)
	}

	return out, nil
}

// ParsePartUUID разбирает строковый UUID из rpc.
func ParsePartUUID(s string) (uuid.UUID, error) {
	if s == "" {
		return uuid.Nil, errs.ErrEmptyUUID
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, errs.ErrInvalidUUID
	}

	return id, nil
}

// PartToProto переводит доменную модель в protobuf Part.
func PartToProto(p model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          p.UUID().String(),
		Name:          p.Name(),
		Description:   p.Description(),
		Price:         p.Price(),
		PartType:      PartTypeToProto(p.PartType()),
		StockQuantity: int64(p.StockQuantity()),
		CreatedAt:     timestamppb.New(p.CreatedAt()),
	}
}

// PartsToProto переводит срез доменных моделей в protobuf.
func PartsToProto(parts []model.Part) []*inventoryv1.Part {
	out := make([]*inventoryv1.Part, 0, len(parts))
	for _, p := range parts {
		out = append(out, PartToProto(p))
	}

	return out
}

// PartTypeFromProto переводит protobuf enum в доменный тип.
func PartTypeFromProto(p inventoryv1.PartType) model.PartType {
	switch p {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return model.PartTypeWeapon
	default:
		return model.PartTypeUnspecified
	}
}

// PartTypeToProto переводит доменный тип в protobuf enum.
func PartTypeToProto(t model.PartType) inventoryv1.PartType {
	switch t {
	case model.PartTypeHull:
		return inventoryv1.PartType_PART_TYPE_HULL
	case model.PartTypeEngine:
		return inventoryv1.PartType_PART_TYPE_ENGINE
	case model.PartTypeShield:
		return inventoryv1.PartType_PART_TYPE_SHIELD
	case model.PartTypeWeapon:
		return inventoryv1.PartType_PART_TYPE_WEAPON
	default:
		return inventoryv1.PartType_PART_TYPE_UNSPECIFIED
	}
}
