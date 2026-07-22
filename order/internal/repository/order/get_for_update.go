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

func (r *Repository) getOrderForUpdate(ctx context.Context, id uuid.UUID) (record.Order, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM orders
		WHERE uuid = $1
		FOR UPDATE`, orderColumns)

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, id)
	if err != nil {
		return record.Order{}, fmt.Errorf("получить заказ for update: %w", err)
	}
	defer rows.Close()

	rec, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[record.Order])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return record.Order{}, errs.ErrOrderNotFound
		}
		return record.Order{}, fmt.Errorf("получить заказ for update: %w", err)
	}

	return rec, nil
}

// GetForUpdate возвращает заказ с блокировкой строки (SELECT ... FOR UPDATE).
// Должен вызываться внутри открытой транзакции.
func (r *Repository) GetForUpdate(ctx context.Context, id uuid.UUID) (model.Order, error) {
	rec, err := r.getOrderForUpdate(ctx, id)
	if err != nil {
		return model.Order{}, err
	}

	items, err := r.getOrderItems(ctx, id)
	if err != nil {
		return model.Order{}, err
	}

	return repoconv.OrderToModel(rec, items), nil
}
