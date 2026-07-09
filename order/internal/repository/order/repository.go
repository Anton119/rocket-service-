package order

import (
	"context"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Anton119/rocket-service-/order/internal/model"
)

// Repository — PostgreSQL-хранилище заказов.
type Repository struct {
	pool      *pgxpool.Pool
	getter    *trmpgx.CtxGetter
	txManager TxManager
}

// New создаёт репозиторий заказов.
func New(pool *pgxpool.Pool, txManager TxManager) *Repository {
	return &Repository{
		pool:      pool,
		getter:    trmpgx.DefaultCtxGetter,
		txManager: txManager,
	}
}

// Create атомарно сохраняет заказ и его позиции.
func (r *Repository) Create(ctx context.Context, order model.Order) error {
	return r.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := r.createOrder(txCtx, order); err != nil {
			return err
		}
		return r.createOrderItems(txCtx, order)
	})
}
