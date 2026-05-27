package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/inventory/internal/model"
)

// GetPart возвращает деталь по UUID.
func (s *Service) GetPart(ctx context.Context, id uuid.UUID) (model.Part, error) {
	part, err := s.repo.Get(ctx, id)
	if err != nil {
		return model.Part{}, fmt.Errorf("получить деталь: %w", err)
	}

	return part, nil
}
