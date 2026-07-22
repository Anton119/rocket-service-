package ship_assembled

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Anton119/rocket-service-/assembly/internal/model"
	"github.com/Anton119/rocket-service-/platform/pkg/kafka"
	eventsv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/events/v1"
)

// Producer отправляет события ShipAssembled в Kafka.
type Producer struct {
	kafkaProducer KafkaProducer
}

// NewProducer создаёт producer событий ShipAssembled.
func NewProducer(kafkaProducer KafkaProducer) *Producer {
	return &Producer{kafkaProducer: kafkaProducer}
}

// Produce сериализует ShipAssembled и отправляет в Kafka с ключом order_uuid.
func (p *Producer) Produce(ctx context.Context, event model.ShipAssembledEvent) error {
	msg := &eventsv1.ShipAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		BuildTimeSec: event.BuildTimeSec,
		AssembledAt:  timestamppb.New(event.AssembledAt),
		UserUuid:     event.UserUUID,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("сериализовать ShipAssembled: %w", err)
	}

	if err := p.kafkaProducer.Send(ctx, &kafka.Message{
		Key:   []byte(event.OrderUUID),
		Value: payload,
	}); err != nil {
		slog.ErrorContext(ctx, "не удалось отправить ShipAssembled", "error", err, "order_uuid", event.OrderUUID)
		return err
	}

	return nil
}
