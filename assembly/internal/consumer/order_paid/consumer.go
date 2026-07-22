package order_paid

import (
	"context"
	"log/slog"
)

// Service потребляет OrderPaid, собирает корабль и публикует ShipAssembled.
type Service struct {
	consumer  Consumer
	assembler Assembler
	producer  ShipAssembledProducer
}

// NewService создаёт consumer-сервис OrderPaid.
func NewService(consumer Consumer, assembler Assembler, producer ShipAssembledProducer) *Service {
	return &Service{
		consumer:  consumer,
		assembler: assembler,
		producer:  producer,
	}
}

// RunConsumer запускает цикл потребления OrderPaid.
func (s *Service) RunConsumer(ctx context.Context) error {
	slog.InfoContext(ctx, "запуск потребителя OrderPaid")
	return s.consumer.Consume(ctx, s.OrderPaidHandler)
}
