package v1

import (
	"context"

	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

// ReserveParts резервирует детали под заказ (reserved += 1).
func (a *API) ReserveParts(ctx context.Context, req *inventoryv1.ReservePartsRequest) (*inventoryv1.ReservePartsResponse, error) {
	uuids, err := parseUUIDs(req.GetUuids())
	if err != nil {
		return nil, err
	}

	if err := a.svc.ReserveParts(ctx, uuids); err != nil {
		return nil, err
	}

	return &inventoryv1.ReservePartsResponse{}, nil
}
