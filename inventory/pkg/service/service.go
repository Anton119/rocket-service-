package service

import (
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"

	invapi "github.com/Anton119/rocket-service-/inventory/internal/api/inventory/v1"
	partrepo "github.com/Anton119/rocket-service-/inventory/internal/repository/part"
	"github.com/Anton119/rocket-service-/inventory/internal/service/domain"
	partsvc "github.com/Anton119/rocket-service-/inventory/internal/service/part"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

// NewInventoryServer собирает gRPC InventoryService с PostgreSQL-репозиторием.
func NewInventoryServer(pool *pgxpool.Pool, txManager *manager.Manager) inventoryv1.InventoryServiceServer {
	repo := partrepo.New(pool, txManager)
	checker := domain.NewCompatibilityChecker()
	svc := partsvc.NewService(repo, txManager, checker)

	return invapi.NewAPI(svc)
}
