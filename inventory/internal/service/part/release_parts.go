package part

import (
	"context"
	"fmt"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/service/input"
)

// ReleaseParts освобождает резерв деталей по списку UUID в одной транзакции.
func (s *Service) ReleaseParts(ctx context.Context, in input.ReleasePartsInput) error {
	return s.txManager.Do(ctx, func(ctx context.Context) error {
		parts, err := s.repo.List(ctx, input.PartFilter{UUIDs: in.UUIDs})
		if err != nil {
			return fmt.Errorf("получить детали: %w", err)
		}

		if len(parts) != len(in.UUIDs) {
			return errs.ErrPartNotFound
		}

		for i := range parts {
			if err = parts[i].Release(); err != nil {
				return fmt.Errorf("освободить %s: %w", parts[i].Name(), err)
			}
		}

		if err = s.repo.UpdateReservedBatch(ctx, parts); err != nil {
			return fmt.Errorf("сохранить освобождение: %w", err)
		}

		return nil
	})
}
