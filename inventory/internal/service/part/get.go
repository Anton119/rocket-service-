package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/inventory/internal/model"
)

// GetPart возвращает деталь по UUID.
func (s *Service) GetPart(ctx context.Context, id uuid.UUID) (model.Part, error) {
	return s.repo.Get(ctx, id)
}
