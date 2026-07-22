package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
)

// UpdateStockReserved сохраняет stock_quantity и reserved для списка деталей.
func (r *Repository) UpdateStockReserved(ctx context.Context, parts []model.Part) error {
	const query = `
		UPDATE parts
		SET stock_quantity = $1,
		    reserved = $2,
		    updated_at = NOW()
		WHERE uuid = $3`

	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	for _, p := range parts {
		id, err := uuid.Parse(p.UUID)
		if err != nil {
			return fmt.Errorf("разобрать uuid: %w", errs.ErrInvalidUUID)
		}
		tag, err := db.Exec(ctx, query, p.StockQuantity, p.Reserved, id)
		if err != nil {
			return fmt.Errorf("обновить деталь %s: %w", p.UUID, err)
		}
		if tag.RowsAffected() == 0 {
			return errs.ErrPartNotFound
		}
	}
	return nil
}
