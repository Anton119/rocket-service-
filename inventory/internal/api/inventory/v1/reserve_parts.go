package v1

import (
	"context"

	apiconv "github.com/Anton119/rocket-service-/inventory/internal/api/converter"
	"github.com/Anton119/rocket-service-/inventory/internal/service/input"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

// ReserveParts реализует rpc ReserveParts.
func (a *API) ReserveParts(
	ctx context.Context,
	req *inventoryv1.ReservePartsRequest,
) (*inventoryv1.ReservePartsResponse, error) {
	ids, err := apiconv.UUIDsFromStrings(req.GetUuids())
	if err != nil {
		return nil, err
	}

	if err = a.svc.ReserveParts(ctx, input.ReservePartsInput{UUIDs: ids}); err != nil {
		return nil, err
	}

	return &inventoryv1.ReservePartsResponse{}, nil
}
