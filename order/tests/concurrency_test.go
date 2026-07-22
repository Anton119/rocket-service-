package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestConcurrency_SelectForUpdateInRepository проверяет, что pay/cancel используют
// блокировку строки через SELECT ... FOR UPDATE в репозитории.
func TestConcurrency_SelectForUpdateInRepository(t *testing.T) {
	t.Helper()

	moduleRoot, err := findModuleRoot()
	if err != nil {
		t.Fatalf("найти корень модуля order: %v", err)
	}

	forUpdatePath := filepath.Join(moduleRoot, "internal/repository/order/get_for_update.go")
	content, readErr := os.ReadFile(forUpdatePath)
	if readErr != nil {
		t.Fatalf("прочитать get_for_update.go: %v", readErr)
	}

	source := string(content)
	if !strings.Contains(source, "FOR UPDATE") {
		t.Fatal("ожидали SELECT ... FOR UPDATE в get_for_update.go")
	}

	payPath := filepath.Join(moduleRoot, "internal/service/order/pay.go")
	paySource, readErr := os.ReadFile(payPath)
	if readErr != nil {
		t.Fatalf("прочитать pay.go: %v", readErr)
	}
	if !strings.Contains(string(paySource), "GetForUpdate") {
		t.Fatal("PayOrder должен вызывать GetForUpdate для конкурентной безопасности")
	}

	cancelPath := filepath.Join(moduleRoot, "internal/service/order/cancel.go")
	cancelSource, readErr := os.ReadFile(cancelPath)
	if readErr != nil {
		t.Fatalf("прочитать cancel.go: %v", readErr)
	}
	if !strings.Contains(string(cancelSource), "GetForUpdate") {
		t.Fatal("CancelOrder должен вызывать GetForUpdate для конкурентной безопасности")
	}
}
