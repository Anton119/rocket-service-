package ship_assembled

import (
	"context"

	"github.com/Anton119/rocket-service-/platform/pkg/kafka"
)

// KafkaProducer — контракт отправки сообщений в Kafka.
type KafkaProducer interface {
	Send(ctx context.Context, msg *kafka.Message) error
}
