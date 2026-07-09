package part

import (
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository — PostgreSQL-хранилище деталей.
type Repository struct {
	pool      *pgxpool.Pool
	getter    *trmpgx.CtxGetter
	txManager TxManager
}

// New создаёт репозиторий деталей.
func New(pool *pgxpool.Pool, txManager TxManager) *Repository {
	return &Repository{
		pool:      pool,
		getter:    trmpgx.DefaultCtxGetter,
		txManager: txManager,
	}
}
