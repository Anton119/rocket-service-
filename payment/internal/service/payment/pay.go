package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/payment/internal/model"
	"github.com/Anton119/rocket-service-/payment/internal/service/input"
)

// PayOrder проводит оплату и возвращает идентификатор транзакции.
func (s *Service) PayOrder(ctx context.Context, in input.PayOrderInput) (uuid.UUID, error) {
	transactionUUID := uuid.New()

	slog.InfoContext(ctx, "оплата заказа",
		slog.String("order_uuid", in.OrderUUID.String()),
		slog.String("payment_method", paymentMethodName(in.Method)),
		slog.String("transaction_uuid", transactionUUID.String()),
	)

	return transactionUUID, nil
}

func paymentMethodName(method model.PaymentMethod) string {
	switch method {
	case model.PaymentMethodCard:
		return "CARD"
	case model.PaymentMethodSBP:
		return "SBP"
	case model.PaymentMethodCreditCard:
		return "CREDIT_CARD"
	case model.PaymentMethodInvestorMoney:
		return "INVESTOR_MONEY"
	default:
		return "UNKNOWN"
	}
}
