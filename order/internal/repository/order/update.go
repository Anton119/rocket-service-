package order

import (
	"context"
	"fmt"
	"time"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
	repoconv "github.com/Anton119/rocket-service-/order/internal/repository/converter"
)

func (r *Repository) updateOrder(ctx context.Context, order model.Order) error {
	rec := repoconv.OrderToRecord(order)
	now := time.Now()

	const query = `
		UPDATE orders
		SET status = $1, transaction_uuid = $2, payment_method = $3, updated_at = $4
		WHERE uuid = $5`

	tag, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query,
		rec.Status,
		rec.TransactionUUID,
		rec.PaymentMethod,
		now,
		rec.UUID,
	)
	if err != nil {
		return fmt.Errorf("обновить заказ: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errs.ErrOrderNotFound
	}

	return nil
}

// Save обновляет заголовок заказа (оплата, отмена).
func (r *Repository) Save(ctx context.Context, order model.Order) error {
	return r.updateOrder(ctx, order)
}
