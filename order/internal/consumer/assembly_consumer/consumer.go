package assembly_consumer

import (
	"context"
	"log/slog"
)

// Service потребляет ShipAssembled и делегирует обработку в OrderService.
type Service struct {
	consumer  Consumer
	assembler OrderAssembler
}

// NewService создаёт consumer-сервис ShipAssembled.
func NewService(consumer Consumer, assembler OrderAssembler) *Service {
	return &Service{
		consumer:  consumer,
		assembler: assembler,
	}
}

// RunConsumer запускает цикл потребления ShipAssembled.
func (s *Service) RunConsumer(ctx context.Context) error {
	slog.InfoContext(ctx, "запуск потребителя ShipAssembled")
	return s.consumer.Consume(ctx, s.ShipAssembledHandler)
}
