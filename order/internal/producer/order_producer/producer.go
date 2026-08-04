package order_producer

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Anton119/rocket-service-/order/internal/model"
	"github.com/Anton119/rocket-service-/platform/pkg/kafka"
	kafkamw "github.com/Anton119/rocket-service-/platform/pkg/middleware/kafka"
	eventsv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/events/v1"
)

// Producer отправляет события OrderPaid в Kafka.
type Producer struct {
	kafkaProducer KafkaProducer
}

// NewProducer создаёт producer OrderPaid.
func NewProducer(kafkaProducer KafkaProducer) *Producer {
	return &Producer{kafkaProducer: kafkaProducer}
}

// Produce сериализует OrderPaid и отправляет с ключом order_uuid.
func (p *Producer) Produce(ctx context.Context, event model.OrderPaidEvent) error {
	msg := &eventsv1.OrderPaid{
		EventUuid:       event.EventUUID,
		OrderUuid:       event.OrderUUID,
		TransactionUuid: event.TransactionUUID,
		PaymentMethod:   event.PaymentMethod,
		PaidAt:          timestamppb.New(event.PaidAt),
		UserUuid:        event.UserUUID,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("сериализовать OrderPaid: %w", err)
	}

	if err := p.kafkaProducer.Send(ctx, &kafka.Message{
		Key:     []byte(event.OrderUUID),
		Value:   payload,
		Headers: kafkamw.ProducerSessionHeaders(ctx),
	}); err != nil {
		slog.ErrorContext(ctx, "не удалось отправить OrderPaid", "error", err, "order_uuid", event.OrderUUID)
		return err
	}

	return nil
}
