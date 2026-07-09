package converter

import (
	"fmt"
	"math"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
	"github.com/Anton119/rocket-service-/inventory/internal/repository/record"
)

func PartModelToRecord(p model.Part) (record.Part, error) {
	id, err := uuid.Parse(p.UUID)
	if err != nil {
		return record.Part{}, fmt.Errorf("разобрать uuid: %w", errs.ErrInvalidUUID)
	}

	if p.StockQuantity > math.MaxInt32 || p.StockQuantity < math.MinInt32 {
		return record.Part{}, fmt.Errorf("stock_quantity вне диапазона int32: %d", p.StockQuantity)
	}

	return record.Part{
		UUID:          id,
		Name:          p.Name,
		Description:   p.Description,
		PartType:      partTypeToString(p.PartType),
		Price:         p.Price,
		StockQuantity: int32(p.StockQuantity),
		CreatedAt:     p.CreatedAt,
	}, nil
}

func PartRecordToModel(p record.Part) model.Part {
	return model.Part{
		UUID:          p.UUID.String(),
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		PartType:      partTypeFromString(p.PartType),
		StockQuantity: int64(p.StockQuantity),
		CreatedAt:     p.CreatedAt,
	}
}

// PartTypeToDB переводит доменный тип детали в значение колонки part_type.
func PartTypeToDB(t model.PartType) string {
	return partTypeToString(t)
}

func partTypeToString(t model.PartType) string {
	switch t {
	case model.PartTypeHull:
		return "HULL"
	case model.PartTypeEngine:
		return "ENGINE"
	case model.PartTypeShield:
		return "SHIELD"
	case model.PartTypeWeapon:
		return "WEAPON"
	default:
		return ""
	}
}

func partTypeFromString(s string) model.PartType {
	switch s {
	case "HULL":
		return model.PartTypeHull
	case "ENGINE":
		return model.PartTypeEngine
	case "SHIELD":
		return model.PartTypeShield
	case "WEAPON":
		return model.PartTypeWeapon
	default:
		return model.PartTypeUnspecified
	}
}
