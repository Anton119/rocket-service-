package converter

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
	"github.com/Anton119/rocket-service-/inventory/internal/repository/record"
)

func PartRecordsToModels(recs []record.Part) ([]model.Part, error) {
	parts := make([]model.Part, 0, len(recs))
	for i := range recs {
		part, err := PartRecordToModel(recs[i])
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}

	return parts, nil
}

func PartRecordToModel(rec record.Part) (model.Part, error) {
	var propsRec record.PartPropertiesRecord
	if len(rec.Properties) > 0 {
		if err := json.Unmarshal(rec.Properties, &propsRec); err != nil {
			return model.Part{}, fmt.Errorf("десериализовать свойства: %w", err)
		}
	}

	props, err := partPropertiesFromRecord(propsRec)
	if err != nil {
		return model.Part{}, fmt.Errorf("конвертировать свойства: %w", err)
	}

	partType, err := model.NewPartType(rec.PartType)
	if err != nil {
		return model.Part{}, fmt.Errorf("конвертировать тип детали: %w", err)
	}

	return model.RestorePart(
		rec.UUID,
		rec.Name,
		rec.Description,
		partType,
		rec.Price,
		int(rec.StockQuantity),
		int(rec.Reserved),
		props,
		rec.CreatedAt,
	), nil
}

func partPropertiesFromRecord(rec record.PartPropertiesRecord) (model.PartProperties, error) {
	switch {
	case rec.Hull != nil:
		return model.NewHullProperties(rec.Hull.Strength)
	case rec.Engine != nil:
		return model.NewEngineProperties(model.EngineClass(rec.Engine.Class), rec.Engine.RequiredStrength)
	case rec.Shield != nil:
		return model.NewShieldProperties(model.ShieldType(rec.Shield.ShieldType))
	case rec.Weapon != nil:
		return model.NewWeaponProperties(model.WeaponType(rec.Weapon.WeaponType))
	default:
		return model.PartProperties{}, nil
	}
}

// PartTypeToDB переводит доменный тип детали в значение колонки part_type.
func PartTypeToDB(t model.PartType) string {
	return string(t)
}

// ParsePartUUID разбирает UUID из строки БД/запроса.
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
