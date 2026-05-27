package v1

import (
	"time"

	orderv1 "github.com/Anton119/rocket-service-/shared/pkg/openapi/order/v1"
)

const (
	inventoryListPartsTimeout = 5 * time.Second
	paymentPayOrderTimeout    = 10 * time.Second
)

// API — HTTP handler (ogen Handler) поверх доменного сервиса.
type API struct {
	orderv1.UnimplementedHandler
	svc OrderService
}

// NewAPI создаёт входной адаптер заказов.
func NewAPI(svc OrderService) *API {
	return &API{svc: svc}
}

// NewServer собирает HTTP-сервер OpenAPI.
func NewServer(api *API) (*orderv1.Server, error) {
	return orderv1.NewServer(api)
}
