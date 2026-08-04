package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/IBM/sarama"

	"github.com/Anton119/rocket-service-/assembly/internal/config"
	orderpaid "github.com/Anton119/rocket-service-/assembly/internal/consumer/order_paid"
	shipassembled "github.com/Anton119/rocket-service-/assembly/internal/producer/ship_assembled"
	assemblysvc "github.com/Anton119/rocket-service-/assembly/internal/service/assembly"
	"github.com/Anton119/rocket-service-/platform/pkg/closer"
	wrappedKafkaConsumer "github.com/Anton119/rocket-service-/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/Anton119/rocket-service-/platform/pkg/kafka/producer"
	kafkaMiddleware "github.com/Anton119/rocket-service-/platform/pkg/middleware/kafka"
)

// ConsumerService — контракт запуска Kafka-потребителя.
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type diContainer struct {
	syncProducer  sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup

	orderPaidKafkaProducer *wrappedKafkaProducer.Producer
	orderPaidKafkaConsumer *wrappedKafkaConsumer.Consumer

	shipAssembledProducer orderpaid.ShipAssembledProducer
	assemblyService       orderpaid.Assembler
	orderPaidConsumer     ConsumerService
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers,
			config.AppConfig().ShipAssembledProducer.SaramaConfig(),
		)
		if err != nil {
			slog.Error("не удалось создать sync producer", "error", err)
			os.Exit(1)
		}
		closer.Add("Kafka sync producer", func(_ context.Context) error {
			return p.Close()
		})
		d.syncProducer = p
	}
	return d.syncProducer
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		cg, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers,
			config.AppConfig().OrderPaidConsumer.ConsumerGroupID(),
			config.AppConfig().OrderPaidConsumer.SaramaConfig(),
		)
		if err != nil {
			slog.Error("не удалось создать consumer group", "error", err)
			os.Exit(1)
		}
		closer.Add("Kafka consumer group", func(_ context.Context) error {
			return cg.Close()
		})
		d.consumerGroup = cg
	}
	return d.consumerGroup
}

func (d *diContainer) ShipAssembledKafkaProducer() *wrappedKafkaProducer.Producer {
	if d.orderPaidKafkaProducer == nil {
		d.orderPaidKafkaProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().ShipAssembledProducer.TopicName(),
		)
	}
	return d.orderPaidKafkaProducer
}

func (d *diContainer) OrderPaidKafkaConsumer() *wrappedKafkaConsumer.Consumer {
	if d.orderPaidKafkaConsumer == nil {
		d.orderPaidKafkaConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{config.AppConfig().OrderPaidConsumer.TopicName()},
			wrappedKafkaConsumer.WithMiddlewares(
				kafkaMiddleware.ConsumerSession(),
				kafkaMiddleware.ConsumerLogging(),
			),
		)
	}
	return d.orderPaidKafkaConsumer
}

func (d *diContainer) ShipAssembledProducer() orderpaid.ShipAssembledProducer {
	if d.shipAssembledProducer == nil {
		d.shipAssembledProducer = shipassembled.NewProducer(d.ShipAssembledKafkaProducer())
	}
	return d.shipAssembledProducer
}

func (d *diContainer) AssemblyService() orderpaid.Assembler {
	if d.assemblyService == nil {
		d.assemblyService = assemblysvc.NewService(&config.AppConfig().Assembler)
	}
	return d.assemblyService
}

func (d *diContainer) OrderPaidConsumerService() ConsumerService {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = orderpaid.NewService(
			d.OrderPaidKafkaConsumer(),
			d.AssemblyService(),
			d.ShipAssembledProducer(),
		)
	}
	return d.orderPaidConsumer
}
