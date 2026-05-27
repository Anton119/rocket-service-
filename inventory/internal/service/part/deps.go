package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/inventory/internal/model"
)

// PartRepository — доступ к каталогу деталей.
type PartRepository interface {
	Get(ctx context.Context, id uuid.UUID) (model.Part, error)
	ListParts(ctx context.Context, partType model.PartType, ids []uuid.UUID) ([]model.Part, error)
}
