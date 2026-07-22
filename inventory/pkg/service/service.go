package service

import (
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"

	invapi "github.com/Anton119/rocket-service-/inventory/internal/api/inventory/v1"
	partrepo "github.com/Anton119/rocket-service-/inventory/internal/repository/part"
	partsvc "github.com/Anton119/rocket-service-/inventory/internal/service/application/part"
	"github.com/Anton119/rocket-service-/inventory/internal/service/domain"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

// NewInventoryServer собирает gRPC InventoryService с PostgreSQL-репозиторием.
func NewInventoryServer(pool *pgxpool.Pool, txManager *manager.Manager) inventoryv1.InventoryServiceServer {
	repo := partrepo.New(pool, txManager)
	svc := partsvc.NewService(repo, txManager, domain.NewCompatibilityChecker())

	return invapi.NewAPI(svc)
}
