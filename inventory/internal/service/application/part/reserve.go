package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
)

// ReserveParts резервирует детали: ListForUpdate → reserved += 1.
func (s *Service) ReserveParts(ctx context.Context, uuids []uuid.UUID) error {
	if len(uuids) == 0 {
		return nil
	}

	return s.tx.Do(ctx, func(txCtx context.Context) error {
		parts, err := s.repo.ListForUpdate(txCtx, uuids)
		if err != nil {
			return err
		}

		for i := range parts {
			if parts[i].Available() < 1 {
				return errs.ErrOutOfStock
			}
			parts[i].Reserved++
		}

		if err := s.repo.UpdateStockReserved(txCtx, parts); err != nil {
			return fmt.Errorf("сохранить резерв: %w", err)
		}
		return nil
	})
}
