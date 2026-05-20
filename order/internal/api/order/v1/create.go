package v1

import (
	"context"
	"errors"
	"net/http"

	apiconv "github.com/Anton119/rocket-service-/order/internal/api/converter"
	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	orderv1 "github.com/Anton119/rocket-service-/shared/pkg/openapi/order/v1"
)

// CreateOrder POST /api/v1/orders.
func (a *API) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	invCtx, cancel := context.WithTimeout(ctx, inventoryListPartsTimeout)
	defer cancel()

	in := apiconv.CreateOrderInputFromRequest(req)
	out, err := a.svc.CreateOrder(invCtx, in)
	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return &orderv1.CreateOrderNotFound{
				Code:    http.StatusNotFound,
				Message: "деталь не найдена",
			}, nil
		}

		var ia *errs.InvalidArgumentError
		if errors.As(err, &ia) {
			return &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: ia.Message,
			}, nil
		}

		if errors.Is(err, errs.ErrPartOutOfStock) {
			return &orderv1.CreateOrderConflict{
				Code:    http.StatusConflict,
				Message: "деталь отсутствует на складе",
			}, nil
		}

		return nil, err
	}

	resp := orderv1.CreateOrderResponse{}
	resp.SetOrderUUID(out.OrderUUID)
	resp.SetTotalPrice(out.TotalPrice)

	return &resp, nil
}
