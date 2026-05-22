package v1

import (
	"context"
	"net/http"

	apiconv "github.com/Anton119/rocket-service-/order/internal/api/converter"
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
		return mapPayOrderError(err)
	}

	out := orderv1.PayOrderResponse{}
	out.SetTransactionUUID(txID)

	return &out, nil
}
