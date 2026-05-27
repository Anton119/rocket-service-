package v1

import (
	"context"

	orderv1 "github.com/Anton119/rocket-service-/shared/pkg/openapi/order/v1"
)

// CancelOrder POST /api/v1/orders/{order_uuid}/cancel.
func (a *API) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	if err := a.svc.CancelOrder(ctx, params.OrderUUID); err != nil {
		return mapCancelOrderError(err)
	}

	return &orderv1.CancelOrderResponse{}, nil
}
