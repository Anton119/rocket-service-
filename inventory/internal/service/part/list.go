package part

import (
	"context"

	"github.com/Anton119/rocket-service-/inventory/internal/model"
	"github.com/Anton119/rocket-service-/inventory/internal/service/input"
)

// ListParts возвращает список деталей по фильтру или упорядоченному списку UUID.
func (s *Service) ListParts(ctx context.Context, in input.ListPartsInput) ([]model.Part, error) {
	return s.repo.ListParts(ctx, in.PartType, in.IDs)
}
