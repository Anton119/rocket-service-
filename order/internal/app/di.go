package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/IBM/sarama"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderapi "github.com/Anton119/rocket-service-/order/internal/api/order/v1"
	authgrpc "github.com/Anton119/rocket-service-/order/internal/client/grpc/auth/v1"
	inventorygrpc "github.com/Anton119/rocket-service-/order/internal/client/grpc/inventory/v1"
	paymentgrpc "github.com/Anton119/rocket-service-/order/internal/client/grpc/payment/v1"
	"github.com/Anton119/rocket-service-/order/internal/config"
	assemblyconsumer "github.com/Anton119/rocket-service-/order/internal/consumer/assembly_consumer"
	orderinterceptor "github.com/Anton119/rocket-service-/order/internal/interceptor"
	ordermiddleware "github.com/Anton119/rocket-service-/order/internal/middleware"
	orderproducer "github.com/Anton119/rocket-service-/order/internal/producer/order_producer"
	orderrepo "github.com/Anton119/rocket-service-/order/internal/repository/order"
	ordersvc "github.com/Anton119/rocket-service-/order/internal/service/order"
	"github.com/Anton119/rocket-service-/platform/pkg/closer"
	wrappedKafkaConsumer "github.com/Anton119/rocket-service-/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/Anton119/rocket-service-/platform/pkg/kafka/producer"
	kafkaMiddleware "github.com/Anton119/rocket-service-/platform/pkg/middleware/kafka"
	authv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/auth/v1"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/payment/v1"
)

// ConsumerService — контракт запуска Kafka-потребителя.
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type diContainer struct {
	pgPool        *pgxpool.Pool
	txManager     *manager.Manager
	syncProducer  sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup

	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn
	iamConn       *grpc.ClientConn

	authClient *authgrpc.Client

	orderPaidKafkaProducer     *wrappedKafkaProducer.Producer
	shipAssembledKafkaConsumer *wrappedKafkaConsumer.Consumer

	orderPaidProducer ordersvc.OrderPaidProducer
	orderService      *ordersvc.Service
	assemblyConsumer  ConsumerService
}

func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		dsn, err := config.AppConfig().PG.DSN()
		if err != nil {
			slog.Error("получение DSN PostgreSQL", "error", err)
			os.Exit(1)
		}

		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			slog.Error("не удалось подключиться к PostgreSQL", "error", err)
			os.Exit(1)
		}

		if err = pool.Ping(ctx); err != nil {
			pool.Close()
			slog.Error("не удалось выполнить ping PostgreSQL", "error", err)
			os.Exit(1)
		}

		closer.Add("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		slog.Info("подключение к PostgreSQL установлено")
		d.pgPool = pool
	}

	return d.pgPool
}

func (d *diContainer) TxManager(ctx context.Context) *manager.Manager {
	if d.txManager == nil {
		txManager, err := manager.New(trmpgx.NewDefaultFactory(d.PGPool(ctx)))
		if err != nil {
			slog.Error("создание transaction manager", "error", err)
			os.Exit(1)
		}
		d.txManager = txManager
	}
	return d.txManager
}

func (d *diContainer) InventoryConn() *grpc.ClientConn {
	if d.inventoryConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryClient.GRPCAddress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(orderinterceptor.SessionForwarder()),
		)
		if err != nil {
			slog.Error("не удалось подключиться к InventoryService", "error", err)
			os.Exit(1)
		}
		closer.Add("Inventory gRPC", func(_ context.Context) error {
			return conn.Close()
		})
		d.inventoryConn = conn
	}
	return d.inventoryConn
}

func (d *diContainer) PaymentConn() *grpc.ClientConn {
	if d.paymentConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentClient.GRPCAddress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			slog.Error("не удалось подключиться к PaymentService", "error", err)
			os.Exit(1)
		}
		closer.Add("Payment gRPC", func(_ context.Context) error {
			return conn.Close()
		})
		d.paymentConn = conn
	}
	return d.paymentConn
}

func (d *diContainer) IAMConn() *grpc.ClientConn {
	if d.iamConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().IAMClient.GRPCAddress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			slog.Error("не удалось подключиться к IAMService", "error", err)
			os.Exit(1)
		}
		closer.Add("IAM gRPC", func(_ context.Context) error {
			return conn.Close()
		})
		d.iamConn = conn
	}
	return d.iamConn
}

func (d *diContainer) AuthClient() *authgrpc.Client {
	if d.authClient == nil {
		d.authClient = authgrpc.NewClient(authv1.NewAuthServiceClient(d.IAMConn()))
	}
	return d.authClient
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers,
			config.AppConfig().OrderPaidProducer.SaramaConfig(),
		)
		if err != nil {
			slog.Error("не удалось создать Kafka sync producer", "error", err)
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
			config.AppConfig().ShipAssembledConsumer.ConsumerGroupID(),
			config.AppConfig().ShipAssembledConsumer.SaramaConfig(),
		)
		if err != nil {
			slog.Error("не удалось создать Kafka consumer group", "error", err)
			os.Exit(1)
		}
		closer.Add("Kafka consumer group", func(_ context.Context) error {
			return cg.Close()
		})
		d.consumerGroup = cg
	}
	return d.consumerGroup
}

func (d *diContainer) OrderPaidKafkaProducer() *wrappedKafkaProducer.Producer {
	if d.orderPaidKafkaProducer == nil {
		d.orderPaidKafkaProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.TopicName(),
		)
	}
	return d.orderPaidKafkaProducer
}

func (d *diContainer) ShipAssembledKafkaConsumer() *wrappedKafkaConsumer.Consumer {
	if d.shipAssembledKafkaConsumer == nil {
		d.shipAssembledKafkaConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{config.AppConfig().ShipAssembledConsumer.TopicName()},
			wrappedKafkaConsumer.WithMiddlewares(
				kafkaMiddleware.ConsumerSession(),
				kafkaMiddleware.ConsumerLogging(),
			),
		)
	}
	return d.shipAssembledKafkaConsumer
}

func (d *diContainer) OrderPaidProducer() ordersvc.OrderPaidProducer {
	if d.orderPaidProducer == nil {
		d.orderPaidProducer = orderproducer.NewProducer(d.OrderPaidKafkaProducer())
	}
	return d.orderPaidProducer
}

func (d *diContainer) OrderService(ctx context.Context) *ordersvc.Service {
	if d.orderService == nil {
		repo := orderrepo.New(d.PGPool(ctx), d.TxManager(ctx))
		inv := inventorygrpc.NewClient(inventoryv1.NewInventoryServiceClient(d.InventoryConn()))
		pay := paymentgrpc.NewClient(paymentv1.NewPaymentServiceClient(d.PaymentConn()))
		d.orderService = ordersvc.NewService(repo, inv, pay, d.OrderPaidProducer(), d.TxManager(ctx))
	}
	return d.orderService
}

func (d *diContainer) AssemblyConsumerService() ConsumerService {
	if d.assemblyConsumer == nil {
		d.assemblyConsumer = assemblyconsumer.NewService(
			d.ShipAssembledKafkaConsumer(),
			d.OrderService(context.Background()),
		)
	}
	return d.assemblyConsumer
}

func (d *diContainer) HTTPServer(ctx context.Context) *http.Server {
	api := orderapi.NewAPI(d.OrderService(ctx))
	orderServer, err := orderapi.NewServer(api)
	if err != nil {
		slog.Error("ошибка создания сервера OpenAPI", "error", err)
		os.Exit(1)
	}

	httpCfg := config.AppConfig().HTTP
	handler := ordermiddleware.AuthMiddleware(d.AuthClient(), orderServer)
	handler = otelhttp.NewHandler(handler, "order-service")
	srv := &http.Server{
		Addr:              httpCfg.Address(),
		Handler:           handler,
		ReadHeaderTimeout: httpCfg.ReadHeaderTimeout(),
		ReadTimeout:       httpCfg.ReadTimeout(),
		WriteTimeout:      httpCfg.WriteTimeout(),
		IdleTimeout:       httpCfg.IdleTimeout(),
	}
	closer.Add("HTTP server", func(shutdownCtx context.Context) error {
		return srv.Shutdown(shutdownCtx)
	})
	return srv
}
