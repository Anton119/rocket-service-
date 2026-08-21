package order_test

import (
	"os"
	"testing"

	ordersvc "github.com/Anton119/rocket-service-/order/internal/service/order"
	"github.com/Anton119/rocket-service-/platform/pkg/metrics"
)

func TestMain(m *testing.M) {
	metrics.Init("order-service-test")
	ordersvc.InitMetrics()

	os.Exit(m.Run())
}
