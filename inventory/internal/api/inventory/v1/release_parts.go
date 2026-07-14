package v1

import (
	"context"

	apiconv "github.com/Anton119/rocket-service-/inventory/internal/api/converter"
	"github.com/Anton119/rocket-service-/inventory/internal/service/input"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

// ReleaseParts реализует rpc ReleaseParts.
func (a *API) ReleaseParts(
	ctx context.Context,
	req *inventoryv1.ReleasePartsRequest,
) (*inventoryv1.ReleasePartsResponse, error) {
	ids, err := apiconv.UUIDsFromStrings(req.GetUuids())
	if err != nil {
		return nil, err
	}

	if err = a.svc.ReleaseParts(ctx, input.ReleasePartsInput{UUIDs: ids}); err != nil {
		return nil, err
	}

	return &inventoryv1.ReleasePartsResponse{}, nil
}
