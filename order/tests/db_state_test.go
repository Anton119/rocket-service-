package tests

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"

	orderrepo "github.com/Anton119/rocket-service-/order/internal/repository/order"
	"github.com/Anton119/rocket-service-/order/pkg/app"
)

// TestDBState_OrderRepositoryGet проверяет, что репозиторий читает заказы из PostgreSQL.
// Пропускается, если БД недоступна (локально без docker-compose).
func TestDBState_OrderRepositoryGet(t *testing.T) {
	if os.Getenv("ORDER_DB_URI") == "" && os.Getenv("DB_URI") == "" {
		t.Skip("ORDER_DB_URI/DB_URI не задан — пропуск проверки состояния БД")
	}

	ctx := context.Background()
	db, err := app.OpenOrderDB(ctx)
	if err != nil {
		t.Skipf("PostgreSQL недоступен: %v", err)
	}
	defer db.Close()

	repo := orderrepo.New(db.Pool, db.TxManager)
	_, err = repo.Get(ctx, uuid.New())
	if err == nil {
		t.Fatal("ожидали ошибку для несуществующего заказа")
	}
}
