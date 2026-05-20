package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apiconv "github.com/Anton119/rocket-service-/inventory/internal/api/converter"
	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

// ListParts реализует rpc ListParts.
func (a *API) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	in, err := apiconv.ListPartsInputFromRequest(req)
	if err != nil {
		if apiconv.IsInvalidUUID(err) {
			return nil, status.Error(codes.InvalidArgument, "неверный формат uuid")
		}

		return nil, err
	}

	items, err := a.svc.ListParts(ctx, in)
	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return nil, status.Error(codes.NotFound, "деталь не найдена")
		}

		return nil, err
	}

	return &inventoryv1.ListPartsResponse{Parts: apiconv.PartsToProto(items)}, nil
}
