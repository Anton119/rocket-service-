package order

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
	"github.com/Anton119/rocket-service-/order/internal/service/input"
)

// PayOrder в одной транзакции: FOR UPDATE → Payment → UPDATE PAID → OrderPaid в Kafka.
func (s *Service) PayOrder(ctx context.Context, in input.PayOrderInput) (uuid.UUID, error) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "order.Pay")
	defer span.End()

	span.SetAttributes(
		attribute.Stringer("order.uuid", in.OrderUUID),
		attribute.String("order.payment_method", in.PaymentMethodStored),
	)

	if in.Method == model.PaymentMethodInvalid {
		recordSpanError(span, errs.ErrInvalidPaymentMethod)
		return uuid.Nil, errs.ErrInvalidPaymentMethod
	}

	var txID uuid.UUID

	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		stored, err := s.repo.GetForUpdate(txCtx, in.OrderUUID)
		if err != nil {
			return err
		}

		switch stored.Status {
		case model.OrderStatusPaid, model.OrderStatusCancelled, model.OrderStatusAssembled:
			return errs.ErrOrderPayNotAllowed
		case model.OrderStatusPendingPayment:
			// ok
		default:
			return errs.ErrOrderPayNotAllowed
		}

		paidTxID, err := s.pay.PayOrder(txCtx, in.OrderUUID.String(), in.Method)
		if err != nil {
			return err
		}
		txID = paidTxID

		stored.Status = model.OrderStatusPaid
		stored.TransactionUUID = &txID
		pm := in.PaymentMethodStored
		stored.PaymentMethod = &pm

		if err := s.repo.Save(txCtx, stored); err != nil {
			return err
		}

		slog.InfoContext(txCtx, "оплата заказа",
			slog.String("order_uuid", stored.OrderUUID.String()),
			slog.String("user_uuid", stored.UserUUID.String()),
			slog.String("payment_method", pm),
			slog.Int64("total_price", stored.TotalPrice),
		)

		ordersPaidTotal.Add(txCtx, 1)
		ordersRevenueTotal.Add(txCtx, stored.TotalPrice)

		return s.orderPaid.Produce(txCtx, model.OrderPaidEvent{
			EventUUID:       uuid.New().String(),
			OrderUUID:       stored.OrderUUID.String(),
			TransactionUUID: txID.String(),
			PaymentMethod:   pm,
			UserUUID:        stored.UserUUID.String(),
			PaidAt:          time.Now().UTC(),
		})
	})
	if err != nil {
		recordSpanError(span, err)
		return uuid.Nil, err
	}

	span.SetAttributes(attribute.Stringer("order.transaction_uuid", txID))
	span.SetStatus(codes.Ok, "")

	return txID, nil
}
