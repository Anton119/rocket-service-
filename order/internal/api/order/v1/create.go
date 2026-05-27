package v1

import (
	"context"

	apiconv "github.com/Anton119/rocket-service-/order/internal/api/converter"
	orderv1 "github.com/Anton119/rocket-service-/shared/pkg/openapi/order/v1"
)

// CreateOrder POST /api/v1/orders.
func (a *API) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	invCtx, cancel := context.WithTimeout(ctx, inventoryListPartsTimeout)
	defer cancel()

	in := apiconv.CreateOrderInputFromRequest(req)
	out, err := a.svc.CreateOrder(invCtx, in)
	if err != nil {
		return mapCreateOrderError(err)
	}

	resp := orderv1.CreateOrderResponse{}
	resp.SetOrderUUID(out.OrderUUID)
	resp.SetTotalPrice(out.TotalPrice)

	return &resp, nil
}
