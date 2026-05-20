package v1

import (
	"context"
	"errors"
	"net/http"

	apiconv "github.com/Anton119/rocket-service-/order/internal/api/converter"
	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	orderv1 "github.com/Anton119/rocket-service-/shared/pkg/openapi/order/v1"
)

// GetOrder GET /api/v1/orders/{order_uuid}.
func (a *API) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := a.svc.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return &orderv1.GetOrderNotFound{
				Code:    http.StatusNotFound,
				Message: "заказ не найден",
			}, nil
		}

		return nil, err
	}

	return apiconv.OrderToDTO(order), nil
}
