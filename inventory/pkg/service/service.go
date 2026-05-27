package service

import (
	invapi "github.com/Anton119/rocket-service-/inventory/internal/api/inventory/v1"
	partrepo "github.com/Anton119/rocket-service-/inventory/internal/repository/part"
	partsvc "github.com/Anton119/rocket-service-/inventory/internal/service/part"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

// NewInventoryServer собирает gRPC InventoryService с seed-каталогом.
func NewInventoryServer() inventoryv1.InventoryServiceServer {
	repo := partrepo.NewRepository(partrepo.SeedParts())
	svc := partsvc.NewService(repo)

	return invapi.NewAPI(svc)
}
