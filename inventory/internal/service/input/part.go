package input

import (
	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/inventory/internal/model"
)

// PartFilter задаёт параметры фильтрации деталей.
type PartFilter struct {
	UUIDs    []uuid.UUID
	PartType model.PartType
}

// ListPartsInput — параметры выборки списка деталей.
type ListPartsInput struct {
	PartType model.PartType
	IDs      []uuid.UUID
}

// ReservePartsInput — UUID деталей для резервирования.
type ReservePartsInput struct {
	UUIDs []uuid.UUID
}

// ReleasePartsInput — UUID деталей для освобождения резерва.
type ReleasePartsInput struct {
	UUIDs []uuid.UUID
}
