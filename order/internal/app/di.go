package app

import (
	"context"
	"log/slog"
	"os"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderapi "github.com/Anton119/rocket-service-/order/internal/api/order/v1"
	inventorygrpc "github.com/Anton119/rocket-service-/order/internal/client/grpc/inventory/v1"
	paymentgrpc "github.com/Anton119/rocket-service-/order/internal/client/grpc/payment/v1"
	"github.com/Anton119/rocket-service-/order/internal/config"
	orderrepo "github.com/Anton119/rocket-service-/order/internal/repository/order"
	ordersvc "github.com/Anton119/rocket-service-/order/internal/service/order"
	"github.com/Anton119/rocket-service-/platform/pkg/closer"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	cfg *config.Config

	pgPool    *pgxpool.Pool
	txManager *manager.Manager

	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn

	orderRepo       ordersvc.OrderRepository
	inventoryClient ordersvc.InventoryClient
	paymentClient   ordersvc.PaymentClient
	orderSvc        orderapi.OrderService
	orderAPI        *orderapi.API
}

func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		pool, err := pgxpool.New(ctx, d.cfg.PG.DSN())
		if err != nil {
			slog.Error("не удалось подключиться к PostgreSQL", "error", err)
			os.Exit(1)
		}

		err = pool.Ping(ctx)
		if err != nil {
			slog.Error("не удалось выполнить ping PostgreSQL", "error", err)
			os.Exit(1)
		}

		closer.Add("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}

	return d.pgPool
}

func (d *diContainer) TxManager(ctx context.Context) *manager.Manager {
	if d.txManager == nil {
		txManager, err := manager.New(trmpgx.NewDefaultFactory(d.PGPool(ctx)))
		if err != nil {
			slog.Error("не удалось создать transaction manager", "error", err)
			os.Exit(1)
		}

		d.txManager = txManager
	}

	return d.txManager
}

func (d *diContainer) InventoryConn(_ context.Context) *grpc.ClientConn {
	if d.inventoryConn == nil {
		conn, err := grpc.NewClient(
			d.cfg.InventoryClient.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			slog.Error("не удалось подключиться к InventoryService", "error", err)
			os.Exit(1)
		}

		closer.Add("Inventory gRPC client", func(_ context.Context) error {
			return conn.Close()
		})

		d.inventoryConn = conn
	}

	return d.inventoryConn
}

func (d *diContainer) PaymentConn(_ context.Context) *grpc.ClientConn {
	if d.paymentConn == nil {
		conn, err := grpc.NewClient(
			d.cfg.PaymentClient.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			slog.Error("не удалось подключиться к PaymentService", "error", err)
			os.Exit(1)
		}

		closer.Add("Payment gRPC client", func(_ context.Context) error {
			return conn.Close()
		})

		d.paymentConn = conn
	}

	return d.paymentConn
}

func (d *diContainer) OrderRepository(ctx context.Context) ordersvc.OrderRepository {
	if d.orderRepo == nil {
		d.orderRepo = orderrepo.New(d.PGPool(ctx), d.TxManager(ctx))
	}

	return d.orderRepo
}

func (d *diContainer) InventoryClient(ctx context.Context) ordersvc.InventoryClient {
	if d.inventoryClient == nil {
		d.inventoryClient = inventorygrpc.NewClient(
			inventoryv1.NewInventoryServiceClient(d.InventoryConn(ctx)),
		)
	}

	return d.inventoryClient
}

func (d *diContainer) PaymentClient(ctx context.Context) ordersvc.PaymentClient {
	if d.paymentClient == nil {
		d.paymentClient = paymentgrpc.NewClient(
			paymentv1.NewPaymentServiceClient(d.PaymentConn(ctx)),
		)
	}

	return d.paymentClient
}

func (d *diContainer) OrderService(ctx context.Context) orderapi.OrderService {
	if d.orderSvc == nil {
		d.orderSvc = ordersvc.NewService(
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
		)
	}

	return d.orderSvc
}

func (d *diContainer) OrderAPI(ctx context.Context) *orderapi.API {
	if d.orderAPI == nil {
		d.orderAPI = orderapi.NewAPI(d.OrderService(ctx))
	}

	return d.orderAPI
}
