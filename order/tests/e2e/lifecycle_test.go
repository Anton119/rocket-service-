//go:build e2e

package e2e_test

import (
	"os"
	"testing"
)

// TestLifecycle_OrderPaidToShipAssembled — сквозной сценарий через Kafka.
// Требует поднятые order-service, assembly-service, Kafka и PostgreSQL.
func TestLifecycle_OrderPaidToShipAssembled(t *testing.T) {
	if os.Getenv("E2E_ENABLED") != "1" {
		t.Skip("E2E_ENABLED=1 не задан — пропуск e2e lifecycle")
	}

	t.Skip("e2e lifecycle: подключить HTTP + Kafka fixtures в CI")
}
