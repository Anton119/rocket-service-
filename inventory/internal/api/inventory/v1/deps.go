package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/inventory/internal/model"
	"github.com/Anton119/rocket-service-/inventory/internal/service/input"
)

// PartService описывает сервисный слой каталога деталей для gRPC API.
type PartService interface {
	GetPart(ctx context.Context, id uuid.UUID) (model.Part, error)
	ListParts(ctx context.Context, in input.ListPartsInput) ([]model.Part, error)
	ValidateCompatibility(ctx context.Context, uuids []uuid.UUID) error
	ReserveParts(ctx context.Context, uuids []uuid.UUID) error
	ReleaseParts(ctx context.Context, uuids []uuid.UUID) error
	CommitParts(ctx context.Context, uuids []uuid.UUID) error
}
