package assembly_consumer

import (
	"context"
	"log/slog"

	"github.com/Anton119/rocket-service-/platform/pkg/kafka"
)

// ShipAssembledHandler обрабатывает одно сообщение.
//
// Ошибка десериализации → nil (commit).
// Ошибка бизнес-обработки → error (не commit).
// Отмена контекста → ctx.Err().
func (s *Service) ShipAssembledHandler(ctx context.Context, msg kafka.Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	event, err := decodeShipAssembled(msg.Value)
	if err != nil {
		slog.ErrorContext(ctx, "не удалось декодировать ShipAssembled, commit offset",
			"error", err,
			"topic", msg.Topic,
			"offset", msg.Offset,
		)
		return nil
	}

	slog.InfoContext(ctx, "получен ShipAssembled",
		"order_uuid", event.OrderUUID,
		"user_uuid", event.UserUUID,
		"build_time_sec", event.BuildTimeSec,
	)

	return s.assembler.HandleShipAssembled(ctx, event)
}
