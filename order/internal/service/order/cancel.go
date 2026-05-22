package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
)

// CancelOrder отменяет заказ в статусе ожидания оплаты.
func (s *Service) CancelOrder(ctx context.Context, id uuid.UUID) error {
	o, err := s.repo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("получить заказ: %w", err)
	}

	switch o.Status {
	case model.OrderStatusPendingPayment:
		// ok
	case model.OrderStatusPaid:
		return errs.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return errs.ErrOrderAlreadyCancelled
	default:
		return errs.ErrOrderCancelNotAllowed
	}

	o.Status = model.OrderStatusCancelled
	if err := s.repo.Save(ctx, o); err != nil {
		return fmt.Errorf("сохранить заказ: %w", err)
	}

	return nil
}
