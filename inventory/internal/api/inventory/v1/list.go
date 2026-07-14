package v1

import (
	"context"

	apiconv "github.com/Anton119/rocket-service-/inventory/internal/api/converter"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

// ListParts реализует rpc ListParts.
func (a *API) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	in, err := apiconv.ListPartsInputFromRequest(req)
	if err != nil {
		return nil, err
	}

	items, err := a.svc.ListParts(ctx, in)
	if err != nil {
		return nil, err
	}

	return &inventoryv1.ListPartsResponse{Parts: apiconv.PartsToProto(items)}, nil
}
