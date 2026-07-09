package order

import (
	"context"
	"fmt"

	"github.com/Anton119/rocket-service-/order/internal/model"
	repoconv "github.com/Anton119/rocket-service-/order/internal/repository/converter"
)

func (r *Repository) createOrder(ctx context.Context, order model.Order) error {
	rec := repoconv.OrderToRecord(order)

	const query = `
		INSERT INTO orders (uuid, status, created_at)
		VALUES ($1, $2, $3)`

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query,
		rec.UUID,
		rec.Status,
		rec.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("создать заказ: %w", err)
	}

	return nil
}

func (r *Repository) createOrderItems(ctx context.Context, order model.Order) error {
	items := repoconv.OrderItemsToRecord(order)

	const query = `
		INSERT INTO order_items (order_uuid, part_uuid, part_type, price)
		VALUES ($1, $2, $3, $4)`

	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	for _, item := range items {
		_, err := db.Exec(ctx, query,
			item.OrderUUID,
			item.PartUUID,
			item.PartType,
			item.Price,
		)
		if err != nil {
			return fmt.Errorf("создать позицию заказа: %w", err)
		}
	}

	return nil
}
