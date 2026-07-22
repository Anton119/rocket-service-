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

// ListParts возвращает детали по списку ids (если ids непустой — порядок как в ids).
// Если ids пустой — выборка по partType и сортировка по имени.
func (r *Repository) ListParts(
	ctx context.Context,
	partType model.PartType,
	ids []uuid.UUID,
) ([]model.Part, error) {
	if len(ids) > 0 {
		return r.listPartsByIDs(ctx, ids)
	}

	return r.listPartsByType(ctx, partType)
}

func (r *Repository) listPartsByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Part, error) {
	query := `SELECT` + partSelectColumns + `
		FROM parts
		WHERE uuid = ANY($1)`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("получить детали: %w", err)
	}
	defer rows.Close()

	recs, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.Part])
	if err != nil {
		return nil, fmt.Errorf("получить детали: %w", err)
	}

	byID := make(map[uuid.UUID]record.Part, len(recs))
	for _, rec := range recs {
		byID[rec.UUID] = rec
	}

	out := make([]model.Part, 0, len(ids))
	for _, id := range ids {
		rec, ok := byID[id]
		if !ok {
			return nil, errs.ErrPartNotFound
		}
		out = append(out, repoconv.PartRecordToModel(rec))
	}

	return out, nil
}

func (r *Repository) listPartsByType(ctx context.Context, partType model.PartType) ([]model.Part, error) {
	var (
		query string
		args  []any
	)

	if partType == model.PartTypeUnspecified {
		query = `SELECT` + partSelectColumns + `
			FROM parts
			ORDER BY name`
	} else {
		query = `SELECT` + partSelectColumns + `
			FROM parts
			WHERE part_type = $1
			ORDER BY name`
		args = append(args, repoconv.PartTypeToDB(partType))
	}

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("получить детали: %w", err)
	}
	defer rows.Close()

	recs, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.Part])
	if err != nil {
		return nil, fmt.Errorf("получить детали: %w", err)
	}

	out := make([]model.Part, 0, len(recs))
	for _, rec := range recs {
		out = append(out, repoconv.PartRecordToModel(rec))
	}

	return out, nil
}
