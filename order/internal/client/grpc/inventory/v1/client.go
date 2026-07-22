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

// CommitParts списывает детали со склада после сборки.
func (c *Client) CommitParts(ctx context.Context, uuids []string) error {
	_, err := c.grpc.CommitParts(ctx, &inventoryv1.CommitPartsRequest{Uuids: uuids})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}
		switch st.Code() {
		case codes.NotFound:
			return errs.ErrPartNotFound
		case codes.FailedPrecondition:
			return &errs.InvalidArgumentError{Message: st.Message()}
		case codes.InvalidArgument:
			return &errs.InvalidArgumentError{Message: st.Message()}
		default:
			return err
		}
	}
	return nil
}

// ReserveParts резервирует детали под заказ.
func (c *Client) ReserveParts(ctx context.Context, uuids []string) error {
	_, err := c.grpc.ReserveParts(ctx, &inventoryv1.ReservePartsRequest{Uuids: uuids})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}
		switch st.Code() {
		case codes.NotFound:
			return errs.ErrPartNotFound
		case codes.ResourceExhausted:
			return errs.ErrPartOutOfStock
		case codes.InvalidArgument:
			return &errs.InvalidArgumentError{Message: st.Message()}
		default:
			return err
		}
	}
	return nil
}

// ReleaseParts снимает резерв деталей (отмена заказа).
func (c *Client) ReleaseParts(ctx context.Context, uuids []string) error {
	_, err := c.grpc.ReleaseParts(ctx, &inventoryv1.ReleasePartsRequest{Uuids: uuids})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}
		switch st.Code() {
		case codes.NotFound:
			return errs.ErrPartNotFound
		case codes.FailedPrecondition:
			return &errs.InvalidArgumentError{Message: st.Message()}
		case codes.InvalidArgument:
			return &errs.InvalidArgumentError{Message: st.Message()}
		default:
			return err
		}
	}
	return nil
}
