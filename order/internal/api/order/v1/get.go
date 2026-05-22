package v1

import (
	"context"

	apiconv "github.com/Anton119/rocket-service-/order/internal/api/converter"
	orderv1 "github.com/Anton119/rocket-service-/shared/pkg/openapi/order/v1"
)

// GetOrder GET /api/v1/orders/{order_uuid}.
func (a *API) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := a.svc.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		return mapGetOrderError(err)
	}

	return apiconv.OrderToDTO(order), nil
}
