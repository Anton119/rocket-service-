package order

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
	"github.com/Anton119/rocket-service-/order/internal/service/input"
)

// PayOrder проверяет статус, вызывает оплату и фиксирует PAID.
func (s *Service) PayOrder(ctx context.Context, in input.PayOrderInput) (uuid.UUID, error) {
	if in.Method == model.PaymentMethodInvalid {
		return uuid.Nil, errs.ErrInvalidPaymentMethod
	}

	stored, err := s.repo.Get(ctx, in.OrderUUID)
	if err != nil {
		return uuid.Nil, err
	}

	switch stored.Status {
	case model.OrderStatusPaid, model.OrderStatusCancelled:
		return uuid.Nil, errs.ErrOrderPayNotAllowed
	case model.OrderStatusPendingPayment:
		// ok
	default:
		return uuid.Nil, errs.ErrOrderPayNotAllowed
	}

	txID, err := s.pay.PayOrder(ctx, in.OrderUUID.String(), in.Method)
	if err != nil {
		return uuid.Nil, err
	}

	stored.Status = model.OrderStatusPaid
	stored.TransactionUUID = &txID
	pm := in.PaymentMethodStored
	stored.PaymentMethod = &pm

	if err := s.repo.Save(ctx, stored); err != nil {
		return uuid.Nil, err
	}

	return txID, nil
}
