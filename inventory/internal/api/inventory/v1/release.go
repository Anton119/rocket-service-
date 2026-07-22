package v1

import (
	"context"

	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

// ReleaseParts освобождает ранее зарезервированные детали (reserved -= 1).
func (a *API) ReleaseParts(ctx context.Context, req *inventoryv1.ReleasePartsRequest) (*inventoryv1.ReleasePartsResponse, error) {
	uuids, err := parseUUIDs(req.GetUuids())
	if err != nil {
		return nil, err
	}

	if err := a.svc.ReleaseParts(ctx, uuids); err != nil {
		return nil, err
	}

	return &inventoryv1.ReleasePartsResponse{}, nil
}
