package order_paid

import (
	"context"

	"github.com/Anton119/rocket-service-/assembly/internal/model"
	"github.com/Anton119/rocket-service-/platform/pkg/kafka"
)

// Consumer — контракт потребления сообщений из Kafka.
type Consumer interface {
	Consume(ctx context.Context, handler kafka.MessageHandler) error
}

// Assembler — application-сервис сборки.
type Assembler interface {
	Assemble(ctx context.Context, paid model.OrderPaidEvent) (model.ShipAssembledEvent, error)
}

// ShipAssembledProducer — контракт отправки ShipAssembled.
type ShipAssembledProducer interface {
	Produce(ctx context.Context, event model.ShipAssembledEvent) error
}
