// Package app содержит вспомогательные функции сборки OrderService для интеграционных тестов.
package app

import (
	"context"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"

	invpkg "github.com/Anton119/rocket-service-/inventory/pkg/service"
	paypkg "github.com/Anton119/rocket-service-/payment/pkg/service"
	"github.com/Anton119/rocket-service-/shared/pkg/config"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/payment/v1"
)

// DB — пул PostgreSQL и Transaction Manager для одного сервиса.
type DB struct {
	Pool      *pgxpool.Pool
	TxManager *manager.Manager
}

// OpenInventoryDB подключается к БД inventory (для API-тестов).
func OpenInventoryDB(ctx context.Context) (*DB, error) {
	dsn, err := config.InventoryDBURI()
	if err != nil {
		return nil, err
	}

	return openDB(ctx, dsn)
}

// OpenOrderDB подключается к БД order (для API-тестов).
func OpenOrderDB(ctx context.Context) (*DB, error) {
	dsn, err := config.OrderDBURI()
	if err != nil {
		return nil, err
	}

	return openDB(ctx, dsn)
}

func openDB(ctx context.Context, dsn string) (*DB, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("создание пула соединений: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("проверка соединения с БД: %w", err)
	}

	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("создание transaction manager: %w", err)
	}

	return &DB{Pool: pool, TxManager: txManager}, nil
}

// Close закрывает пул соединений.
func (db *DB) Close() {
	if db != nil && db.Pool != nil {
		db.Pool.Close()
	}
}

// NewInventoryServer собирает gRPC InventoryService (order/tests не импортирует inventory/internal).
func NewInventoryServer(db *DB) inventoryv1.InventoryServiceServer {
	return invpkg.NewInventoryServer(db.Pool, db.TxManager)
}

// NewPaymentServer собирает gRPC PaymentService.
func NewPaymentServer() paymentv1.PaymentServiceServer {
	return paypkg.NewPaymentServer()
}
