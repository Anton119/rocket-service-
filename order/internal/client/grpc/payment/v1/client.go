package paymentgrpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
	paymentv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/payment/v1"
)

// Client — gRPC-клиент оплаты.
type Client struct {
	grpc paymentv1.PaymentServiceClient
}

// NewClient создаёт обёртку над payment PaymentServiceClient.
func NewClient(c paymentv1.PaymentServiceClient) *Client {
	return &Client{grpc: c}
}

func toProtoPaymentMethod(m model.PaymentMethod) (paymentv1.PaymentMethod, bool) {
	switch m {
	case model.PaymentMethodCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD, true
	case model.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP, true
	case model.PaymentMethodCreditCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, true
	case model.PaymentMethodInvestorMoney:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, true
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, false
	}
}

// PayOrder вызывает payment-сервис и возвращает transaction UUID.
func (c *Client) PayOrder(ctx context.Context, orderUUID string, method model.PaymentMethod) (uuid.UUID, error) {
	pm, ok := toProtoPaymentMethod(method)
	if !ok {
		return uuid.Nil, errs.ErrInvalidPaymentMethod
	}

	resp, err := c.grpc.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		PaymentMethod: pm,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.InvalidArgument {
			return uuid.Nil, &errs.InvalidArgumentError{Message: st.Message()}
		}

		return uuid.Nil, err
	}

	txID, err := uuid.Parse(resp.GetTransactionUuid())
	if err != nil {
		return uuid.Nil, err
	}

	return txID, nil
}
