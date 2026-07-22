package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/inventory/internal/model"
)

// TxManager — управление транзакциями БД.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

// PartRepository — доступ к каталогу деталей.
type PartRepository interface {
	Get(ctx context.Context, id uuid.UUID) (model.Part, error)
	ListParts(ctx context.Context, partType model.PartType, ids []uuid.UUID) ([]model.Part, error)
	ListForUpdate(ctx context.Context, ids []uuid.UUID) ([]model.Part, error)
	UpdateStockReserved(ctx context.Context, parts []model.Part) error
}

// CompatibilityChecker — доменные правила совместимости деталей.
type CompatibilityChecker interface {
	Check(slots model.ResolvedShipSlots) error
}
