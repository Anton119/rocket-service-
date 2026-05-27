package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/order/internal/model"
)

// GetOrder возвращает заказ по идентификатору.
func (s *Service) GetOrder(ctx context.Context, id uuid.UUID) (model.Order, error) {
	return s.repo.Get(ctx, id)
}
