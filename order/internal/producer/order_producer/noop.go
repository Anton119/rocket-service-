package order_producer

import (
	"context"

	"github.com/Anton119/rocket-service-/order/internal/model"
)

// NoopProducer — no-op реализация для тестов без Kafka.
type NoopProducer struct{}

// Produce ничего не делает.
func (NoopProducer) Produce(_ context.Context, _ model.OrderPaidEvent) error {
	return nil
}
