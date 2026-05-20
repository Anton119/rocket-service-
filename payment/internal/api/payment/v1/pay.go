package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apiconv "github.com/Anton119/rocket-service-/payment/internal/api/converter"
	paymentv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/payment/v1"
)

// PayOrder реализует rpc PayOrder.
func (a *API) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	in, err := apiconv.PayOrderInputFromRequest(req)
	if err != nil {
		if apiconv.IsValidationError(err) {
			return nil, status.Error(codes.InvalidArgument, apiconv.ValidationErrorMessage(err))
		}

		return nil, err
	}

	out, err := a.svc.PayOrder(ctx, in)
	if err != nil {
		return nil, err
	}

	return apiconv.PayOrderResponseFromResult(*out), nil
}
