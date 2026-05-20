package v1

import (
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

// API — gRPC-адаптер InventoryService поверх доменного сервиса.
type API struct {
	inventoryv1.UnimplementedInventoryServiceServer
	svc PartService
}

// NewAPI создаёт обработчик gRPC.
func NewAPI(svc PartService) *API {
	return &API{svc: svc}
}
