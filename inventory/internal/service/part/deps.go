package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/inventory/internal/model"
	"github.com/Anton119/rocket-service-/inventory/internal/service/input"
)

// PartRepository — доступ к каталогу деталей.
type PartRepository interface {
	Get(ctx context.Context, id uuid.UUID) (model.Part, error)
	List(ctx context.Context, filter input.PartFilter) ([]model.Part, error)
	UpdateReservedBatch(ctx context.Context, parts []model.Part) error
}

// TxManager — выполнение операций в транзакции.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

// CompatibilityChecker — доменный сервис проверки совместимости.
type CompatibilityChecker interface {
	Check(slots model.ResolvedShipSlots) error
}
