package v1

import (
	"context"

	apiconv "github.com/Anton119/rocket-service-/inventory/internal/api/converter"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

// ValidateCompatibility реализует rpc ValidateCompatibility.
func (a *API) ValidateCompatibility(
	ctx context.Context,
	req *inventoryv1.ValidateCompatibilityRequest,
) (*inventoryv1.ValidateCompatibilityResponse, error) {
	slots, err := apiconv.ShipSlotsFromRequest(req)
	if err != nil {
		return nil, err
	}

	if err = a.svc.ValidateCompatibility(ctx, slots); err != nil {
		return nil, err
	}

	return &inventoryv1.ValidateCompatibilityResponse{}, nil
}
