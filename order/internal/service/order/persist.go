package order

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/order/internal/model"
)

func (s *Service) persistOrder(ctx context.Context, o model.Order, partUUIDs []string, userUUID uuid.UUID) error {
	if err := s.repo.Create(ctx, o); err != nil {
		slog.ErrorContext(ctx, "не удалось создать заказ",
			slog.String("error", err.Error()),
			slog.String("user_uuid", userUUID.String()),
		)
		if releaseErr := s.inv.ReleaseParts(ctx, partUUIDs); releaseErr != nil {
			return errors.Join(
				fmt.Errorf("сохранение заказа: %w", err),
				fmt.Errorf("освободить резерв: %w", releaseErr),
			)
		}

		return fmt.Errorf("сохранение заказа: %w", err)
	}

	return nil
}
