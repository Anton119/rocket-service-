package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
	repoconv "github.com/Anton119/rocket-service-/inventory/internal/repository/converter"
	"github.com/Anton119/rocket-service-/inventory/internal/repository/record"
)

// ListForUpdate возвращает детали с блокировкой строк (SELECT ... FOR UPDATE).
// ORDER BY uuid — единый порядок блокировок, чтобы избежать deadlock.
func (r *Repository) ListForUpdate(ctx context.Context, ids []uuid.UUID) ([]model.Part, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	const query = `SELECT` + partSelectColumns + `
		FROM parts
		WHERE uuid = ANY($1)
		ORDER BY uuid
		FOR UPDATE`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("list for update: %w", err)
	}
	defer rows.Close()

	recs, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.Part])
	if err != nil {
		return nil, fmt.Errorf("list for update: %w", err)
	}

	byID := make(map[uuid.UUID]record.Part, len(recs))
	for _, rec := range recs {
		byID[rec.UUID] = rec
	}

	for _, id := range ids {
		if _, ok := byID[id]; !ok {
			return nil, errs.ErrPartNotFound
		}
	}

	out := make([]model.Part, 0, len(recs))
	for _, rec := range recs {
		out = append(out, repoconv.PartRecordToModel(rec))
	}
	return out, nil
}
