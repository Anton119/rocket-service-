package part

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
	repoconv "github.com/Anton119/rocket-service-/inventory/internal/repository/converter"
	"github.com/Anton119/rocket-service-/inventory/internal/repository/record"
)

// Get возвращает деталь по UUID или errs.ErrPartNotFound.
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (model.Part, error) {
	query := `SELECT` + partSelectColumns + `
		FROM parts
		WHERE uuid = $1`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, id)
	if err != nil {
		return model.Part{}, fmt.Errorf("получить деталь: %w", err)
	}
	defer rows.Close()

	rec, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[record.Part])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Part{}, errs.ErrPartNotFound
		}
		return model.Part{}, fmt.Errorf("получить деталь: %w", err)
	}

	return repoconv.PartRecordToModel(rec), nil
}
