package input

import (
	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/inventory/internal/model"
)

// ListPartsInput — параметры выборки списка деталей.
type ListPartsInput struct {
	PartType model.PartType
	IDs      []uuid.UUID
}
