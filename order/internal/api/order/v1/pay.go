package v1

import (
	"context"
	"errors"
	"net/http"

	apiconv "github.com/Anton119/rocket-service-/order/internal/api/converter"
	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	orderv1 "github.com/Anton119/rocket-service-/shared/pkg/openapi/order/v1"
)

// PayOrder POST /api/v1/orders/{order_uuid}/pay.
func (a *API) PayOrder(
	ctx context.Context,
	req *orderv1.PayOrderRequest,
	params orderv1.PayOrderParams,
) (orderv1.PayOrderRes, error) {
	in, ok := apiconv.PayOrderInputFromRequest(req, params.OrderUUID)
	if !ok {
		return &orderv1.PayOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "неизвестный способ оплаты",
		}, nil
	}

	payCtx, cancel := context.WithTimeout(ctx, paymentPayOrderTimeout)
	defer cancel()

	txID, err := a.svc.PayOrder(payCtx, in)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return &orderv1.PayOrderNotFound{
				Code:    http.StatusNotFound,
				Message: "заказ не найден",
			}, nil
		}

		if errors.Is(err, errs.ErrOrderPayNotAllowed) {
			return &orderv1.PayOrderConflict{
				Code:    http.StatusConflict,
				Message: "оплата невозможна в текущем статусе",
			}, nil
		}

		if errors.Is(err, errs.ErrInvalidPaymentMethod) {
			return &orderv1.PayOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: "неизвестный способ оплаты",
			}, nil
		}

		var ia *errs.InvalidArgumentError
		if errors.As(err, &ia) {
			return &orderv1.PayOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: ia.Message,
			}, nil
		}

		return nil, err
	}

	out := orderv1.PayOrderResponse{}
	out.SetTransactionUUID(txID)

	return &out, nil
}
