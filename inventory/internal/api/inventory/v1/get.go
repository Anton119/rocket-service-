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

// GetPart реализует rpc GetPart.
func (a *API) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	id, err := apiconv.ParsePartUUID(req.GetUuid())
	if err != nil {
		if errs.IsInvalidUUID(err) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, err
	}

	part, err := a.svc.GetPart(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return nil, status.Error(codes.NotFound, errs.ErrPartNotFound.Error())
		}

		return nil, err
	}

	return &inventoryv1.GetPartResponse{
		Part: apiconv.PartToProto(part),
	}, nil
}
