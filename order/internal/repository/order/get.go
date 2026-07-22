package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
	repoconv "github.com/Anton119/rocket-service-/order/internal/repository/converter"
	"github.com/Anton119/rocket-service-/order/internal/repository/record"
)

func (r *Repository) getOrder(ctx context.Context, id uuid.UUID) (record.Order, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM orders
		WHERE uuid = $1`, orderColumns)

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, id)
	if err != nil {
		return record.Order{}, fmt.Errorf("получить заказ: %w", err)
	}
	defer rows.Close()

	rec, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[record.Order])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return record.Order{}, errs.ErrOrderNotFound
		}
		return record.Order{}, fmt.Errorf("получить заказ: %w", err)
	}

	return rec, nil
}

func (r *Repository) getOrderItems(ctx context.Context, orderUUID uuid.UUID) ([]record.OrderItem, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM order_items
		WHERE order_uuid = $1`, itemColumns)

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, orderUUID)
	if err != nil {
		return nil, fmt.Errorf("получить позиции заказа: %w", err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.OrderItem])
	if err != nil {
		return nil, fmt.Errorf("получить позиции заказа: %w", err)
	}

	return items, nil
}

// Get возвращает заказ с позициями или errs.ErrOrderNotFound.
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (model.Order, error) {
	rec, err := r.getOrder(ctx, id)
	if err != nil {
		return model.Order{}, err
	}

	items, err := r.getOrderItems(ctx, id)
	if err != nil {
		return model.Order{}, err
	}

	return repoconv.OrderToModel(rec, items), nil
}
