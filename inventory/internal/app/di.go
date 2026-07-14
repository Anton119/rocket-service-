package app

import (
	"context"
	"log/slog"
	"os"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Anton119/rocket-service-/inventory/internal/api/inventory/v1"
	"github.com/Anton119/rocket-service-/inventory/internal/config"
	partrepo "github.com/Anton119/rocket-service-/inventory/internal/repository/part"
	"github.com/Anton119/rocket-service-/inventory/internal/service/domain"
	partsvc "github.com/Anton119/rocket-service-/inventory/internal/service/part"
	"github.com/Anton119/rocket-service-/platform/pkg/closer"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	cfg *config.Config

	pgPool    *pgxpool.Pool
	txManager *manager.Manager

	partRepo partsvc.PartRepository
	partSvc  v1.PartService

	inventoryHandler inventoryv1.InventoryServiceServer
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

func (d *diContainer) PartRepository(ctx context.Context) partsvc.PartRepository {
	if d.partRepo == nil {
		d.partRepo = partrepo.New(d.PGPool(ctx), d.TxManager(ctx))
	}

	return d.partRepo
}

func (d *diContainer) PartService(ctx context.Context) v1.PartService {
	if d.partSvc == nil {
		d.partSvc = partsvc.NewService(
			d.PartRepository(ctx),
			d.TxManager(ctx),
			domain.NewCompatibilityChecker(),
		)
	}

	return d.partSvc
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryv1.InventoryServiceServer {
	if d.inventoryHandler == nil {
		d.inventoryHandler = v1.NewAPI(d.PartService(ctx))
	}

	return d.inventoryHandler
}
