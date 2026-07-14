package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/inventory/internal/model"
)

// UpdateReservedBatch обновляет reserved для нескольких деталей одним запросом.
func (r *Repository) UpdateReservedBatch(ctx context.Context, parts []model.Part) error {
	const query = `
		UPDATE parts AS p
		SET reserved   = batch.reserved,
		    updated_at = NOW()
		FROM unnest($1::uuid[], $2::int[]) AS batch(uuid, reserved)
		WHERE p.uuid = batch.uuid
	`

	uuids := make([]uuid.UUID, len(parts))
	reservedVals := make([]int32, len(parts))

	for i, p := range parts {
		uuids[i] = p.UUID()
		reservedVals[i] = int32(p.Reserved()) //nolint:gosec // G115: reserved ограничен stock_quantity в доменной модели
	}

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, uuids, reservedVals)

	return err
}
