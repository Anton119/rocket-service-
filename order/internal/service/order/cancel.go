package order

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
)

// CancelOrder отменяет заказ в статусе ожидания оплаты.
func (s *Service) CancelOrder(ctx context.Context, id uuid.UUID) error {
	o, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if o.Status != model.OrderStatusPendingPayment {
		return errs.ErrOrderCancelNotAllowed
	}

	o.Status = model.OrderStatusCancelled
	return s.repo.Save(ctx, o)
}
