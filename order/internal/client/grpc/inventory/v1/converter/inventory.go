package converter

import (
	"github.com/Anton119/rocket-service-/order/internal/model"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

func PartProtoToModel(p *inventoryv1.Part) model.Part {
	if p == nil {
		return model.Part{}
	}

	return model.Part{
		UUID:          p.GetUuid(),
		Price:         p.GetPrice(),
		StockQuantity: p.GetStockQuantity(),
	}
}
