package app

import (
	"context"
	"log/slog"
	"os"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	invapi "github.com/Anton119/rocket-service-/inventory/internal/api/inventory/v1"
	authgrpc "github.com/Anton119/rocket-service-/inventory/internal/client/grpc/auth/v1"
	"github.com/Anton119/rocket-service-/inventory/internal/config"
	partrepo "github.com/Anton119/rocket-service-/inventory/internal/repository/part"
	partsvc "github.com/Anton119/rocket-service-/inventory/internal/service/application/part"
	"github.com/Anton119/rocket-service-/inventory/internal/service/domain"
	"github.com/Anton119/rocket-service-/platform/pkg/closer"
	authv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/auth/v1"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

func newGRPCServer(opts ...grpc.ServerOption) *grpc.Server {
	all := append([]grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}, opts...)

	return grpc.NewServer(all...)
}

type diContainer struct {
	pgPool       *pgxpool.Pool
	txManager    *manager.Manager
	iamConn      *grpc.ClientConn
	authClient   *authgrpc.Client
	partRepo     partsvc.PartRepository
	partService  invapi.PartService
	inventoryAPI inventoryv1.InventoryServiceServer
}

func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		dsn, err := config.AppConfig().PG.DSN()
		if err != nil {
			slog.Error("получение DSN", "error", err)
			os.Exit(1)
		}

		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			slog.Error("создание пула соединений", "error", err)
			os.Exit(1)
		}

		if err = pool.Ping(ctx); err != nil {
			pool.Close()
			slog.Error("проверка соединения с БД", "error", err)
			os.Exit(1)
		}
		slog.Info("подключение к PostgreSQL установлено")

		closer.Add("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}
	return d.pgPool
}

func (d *diContainer) TxManager() *manager.Manager {
	if d.txManager == nil {
		txManager, err := manager.New(trmpgx.NewDefaultFactory(d.pgPool))
		if err != nil {
			slog.Error("создание transaction manager", "error", err)
			os.Exit(1)
		}
		d.txManager = txManager
	}
	return d.txManager
}

func (d *diContainer) IAMConn() *grpc.ClientConn {
	if d.iamConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().IAMClient.GRPCAddress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
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

func (d *diContainer) PartRepository(ctx context.Context) partsvc.PartRepository {
	if d.partRepo == nil {
		d.partRepo = partrepo.New(d.PGPool(ctx), d.TxManager())
	}
	return d.partRepo
}

func (d *diContainer) PartService(ctx context.Context) invapi.PartService {
	if d.partService == nil {
		d.partService = partsvc.NewService(
			d.PartRepository(ctx),
			d.TxManager(),
			domain.NewCompatibilityChecker(),
		)
	}
	return d.partService
}

func (d *diContainer) InventoryAPI(ctx context.Context) inventoryv1.InventoryServiceServer {
	if d.inventoryAPI == nil {
		d.inventoryAPI = invapi.NewAPI(d.PartService(ctx))
	}
	return d.inventoryAPI
}
