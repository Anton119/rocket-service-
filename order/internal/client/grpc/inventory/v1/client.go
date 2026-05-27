package inventorygrpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Anton119/rocket-service-/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

// Client — gRPC-клиент каталога деталей.
type Client struct {
	grpc inventoryv1.InventoryServiceClient
}

// NewClient создаёт обёртку над inventory InventoryServiceClient.
func NewClient(c inventoryv1.InventoryServiceClient) *Client {
	return &Client{grpc: c}
}

// ListParts возвращает детали по UUID (доменные модели).
func (c *Client) ListParts(ctx context.Context, uuids []string) ([]model.Part, error) {
	resp, err := c.grpc.ListParts(ctx, &inventoryv1.ListPartsRequest{Uuids: uuids})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, err
		}
		switch st.Code() {
		case codes.NotFound:
			return nil, errs.ErrPartNotFound
		case codes.InvalidArgument:
			return nil, &errs.InvalidArgumentError{Message: st.Message()}
		default:
			return nil, err
		}
	}

	parts := make([]model.Part, 0, len(resp.GetParts()))
	for _, p := range resp.GetParts() {
		parts = append(parts, converter.PartProtoToModel(p))
	}

	return parts, nil
}
