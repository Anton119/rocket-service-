package order

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/order/internal/model"
)

// HandleShipAssembled обрабатывает событие сборки: CommitParts + статус ASSEMBLED.
func (s *Service) HandleShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	orderUUID, err := uuid.Parse(event.OrderUUID)
	if err != nil {
		slog.ErrorContext(ctx, "невалидный order_uuid в ShipAssembled, skip",
			"order_uuid", event.OrderUUID,
			"error", err,
		)
		return nil
	}

	slog.InfoContext(ctx, "обработка ShipAssembled",
		"order_uuid", event.OrderUUID,
		"user_uuid", event.UserUUID,
	)

	order, err := s.repo.Get(ctx, orderUUID)
	if err != nil {
		return err
	}

	if order.Status == model.OrderStatusAssembled {
		slog.InfoContext(ctx, "заказ уже ASSEMBLED, идемпотентный skip", "order_uuid", event.OrderUUID)
		return nil
	}

	if order.Status != model.OrderStatusPaid {
		slog.WarnContext(ctx, "заказ не в статусе PAID, skip ShipAssembled",
			"order_uuid", event.OrderUUID,
			"status", order.Status,
		)
		return nil
	}

	partUUIDs := make([]string, 0, len(order.Items))
	for _, item := range order.Items {
		partUUIDs = append(partUUIDs, item.PartUUID.String())
	}

	if err := s.inv.CommitParts(ctx, partUUIDs); err != nil {
		return fmt.Errorf("CommitParts: %w", err)
	}

	order.Status = model.OrderStatusAssembled
	if err := s.repo.Save(ctx, order); err != nil {
		return fmt.Errorf("сохранить ASSEMBLED: %w", err)
	}

	slog.InfoContext(ctx, "заказ переведён в ASSEMBLED", "order_uuid", event.OrderUUID)
	return nil
}
