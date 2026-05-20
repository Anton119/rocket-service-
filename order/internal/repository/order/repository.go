package order

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
	repoconv "github.com/Anton119/rocket-service-/order/internal/repository/converter"
	"github.com/Anton119/rocket-service-/order/internal/repository/record"
)

// Repository — потокобезопасное in-memory хранилище заказов.
type Repository struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]record.Order
}

// NewRepository создаёт пустое хранилище.
func NewRepository() *Repository {
	return &Repository{
		orders: make(map[uuid.UUID]record.Order),
	}
}

// Get возвращает заказ или errs.ErrOrderNotFound.
func (r *Repository) Get(_ context.Context, id uuid.UUID) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rec, ok := r.orders[id]
	if !ok {
		return model.Order{}, errs.ErrOrderNotFound
	}

	return repoconv.OrderToModel(rec), nil
}

// Save создаёт или обновляет заказ.
func (r *Repository) Save(_ context.Context, o model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[o.OrderUUID] = repoconv.OrderToRecord(o)
	return nil
}
