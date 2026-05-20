package v1

import (
	"context"
	"errors"
	"net/http"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	orderv1 "github.com/Anton119/rocket-service-/shared/pkg/openapi/order/v1"
)

// CancelOrder POST /api/v1/orders/{order_uuid}/cancel.
func (a *API) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	err := a.svc.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return &orderv1.CancelOrderNotFound{
				Code:    http.StatusNotFound,
				Message: "заказ не найден",
			}, nil
		}

		if errors.Is(err, errs.ErrOrderCancelNotAllowed) {
			return &orderv1.CancelOrderConflict{
				Code:    http.StatusConflict,
				Message: "отмена невозможна в текущем статусе",
			}, nil
		}

		return nil, err
	}

	return &orderv1.CancelOrderResponse{}, nil
}
