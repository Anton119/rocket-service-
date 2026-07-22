package v1

import (
	"context"

	apiconv "github.com/Anton119/rocket-service-/inventory/internal/api/converter"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

// GetPart реализует rpc GetPart.
func (a *API) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	id, err := apiconv.ParsePartUUID(req.GetUuid())
	if err != nil {
		return nil, err
	}

	part, err := a.svc.GetPart(ctx, id)
	if err != nil {
		return nil, err
	}

	return &inventoryv1.GetPartResponse{
		Part: apiconv.PartToProto(part),
	}, nil
}
