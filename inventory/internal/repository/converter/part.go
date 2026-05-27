package converter

import (
	"github.com/Anton119/rocket-service-/inventory/internal/model"
	"github.com/Anton119/rocket-service-/inventory/internal/repository/record"
)

func PartModelToRecord(p model.Part) record.Part {
	return record.Part{
		UUID:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		PartType:      int32(p.PartType),
		StockQuantity: p.StockQuantity,
		CreatedAt:     p.CreatedAt,
	}
}

func PartRecordToModel(p record.Part) model.Part {
	return model.Part{
		UUID:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		PartType:      model.PartType(p.PartType),
		StockQuantity: p.StockQuantity,
		CreatedAt:     p.CreatedAt,
	}
}
