package assembly_consumer

import (
	"context"

	"github.com/Anton119/rocket-service-/order/internal/model"
	"github.com/Anton119/rocket-service-/platform/pkg/kafka"
)

// Consumer — контракт потребления сообщений из Kafka.
type Consumer interface {
	Consume(ctx context.Context, handler kafka.MessageHandler) error
}

// OrderAssembler — бизнес-обработка ShipAssembled.
type OrderAssembler interface {
	HandleShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}
