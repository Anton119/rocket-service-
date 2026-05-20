package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/payment/internal/service/input"
)

// PayOrder проводит оплату и возвращает идентификатор транзакции.
func (s *Service) PayOrder(ctx context.Context, in input.PayOrderInput) (*input.PayOrderResult, error) {
	transactionUUID := uuid.New()

	slog.InfoContext(ctx, "оплата прошла успешно",
		"order_uuid", in.OrderUUID.String(),
		"transaction_uuid", transactionUUID.String(),
	)

	return &input.PayOrderResult{
		TransactionUUID: transactionUUID,
	}, nil
}
