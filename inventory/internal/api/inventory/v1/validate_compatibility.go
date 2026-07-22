package v1

import (
	"context"

	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

// ValidateCompatibility проверяет совместимость деталей между собой.
func (a *API) ValidateCompatibility(
	ctx context.Context,
	req *inventoryv1.ValidateCompatibilityRequest,
) (*inventoryv1.ValidateCompatibilityResponse, error) {
	uuids, err := parseUUIDs(req.GetUuids())
	if err != nil {
		return nil, err
	}

	if err := a.svc.ValidateCompatibility(ctx, uuids); err != nil {
		return nil, err
	}

	return &inventoryv1.ValidateCompatibilityResponse{}, nil
}
