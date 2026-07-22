package v1

import (
	"context"

	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

// CommitParts списывает детали со склада (stock -= 1, reserved -= 1).
func (a *API) CommitParts(ctx context.Context, req *inventoryv1.CommitPartsRequest) (*inventoryv1.CommitPartsResponse, error) {
	uuids, err := parseUUIDs(req.GetUuids())
	if err != nil {
		return nil, err
	}

	if err := a.svc.CommitParts(ctx, uuids); err != nil {
		return nil, err
	}

	return &inventoryv1.CommitPartsResponse{}, nil
}
