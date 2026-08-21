package order_paid

import (
	"context"
	"log/slog"

	"github.com/Anton119/rocket-service-/platform/pkg/kafka"
)

// OrderPaidHandler обрабатывает одно сообщение OrderPaid.
//
// Ошибка десериализации → nil (commit).
// Ошибка сборки/отправки → error (не commit).
// Отмена контекста → ctx.Err() (не commit).
func (s *Service) OrderPaidHandler(ctx context.Context, msg kafka.Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	event, err := decodeOrderPaid(msg.Value)
	if err != nil {
		slog.ErrorContext(ctx, "не удалось декодировать OrderPaid, commit offset",
			"error", err,
			"topic", msg.Topic,
			"offset", msg.Offset,
		)
		return nil
	}

	slog.InfoContext(ctx, "начинаем сборку корабля",
		slog.String("order_uuid", event.OrderUUID),
		slog.Int64("offset", msg.Offset),
		slog.Int("partition", int(msg.Partition)),
	)

	assembled, err := s.assembler.Assemble(ctx, event)
	if err != nil {
		return err
	}

	if err := s.producer.Produce(ctx, assembled); err != nil {
		return err
	}

	slog.InfoContext(ctx, "сборка завершена, ShipAssembled отправлен",
		"order_uuid", assembled.OrderUUID,
		"user_uuid", assembled.UserUUID,
		"build_time_sec", assembled.BuildTimeSec,
	)

	return nil
}
