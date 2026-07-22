//go:build e2e

package e2e_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("E2E_ENABLED") != "1" {
		os.Exit(0)
	}
	os.Exit(m.Run())
}
