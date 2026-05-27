package part

import (
	"context"
	"fmt"

	"github.com/Anton119/rocket-service-/inventory/internal/model"
	"github.com/Anton119/rocket-service-/inventory/internal/service/input"
)

// ListParts возвращает список деталей по фильтру или упорядоченному списку UUID.
func (s *Service) ListParts(ctx context.Context, in input.ListPartsInput) ([]model.Part, error) {
	parts, err := s.repo.ListParts(ctx, in.PartType, in.IDs)
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}

	return parts, nil
}
