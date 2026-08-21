package order

import (
	"context"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	metricsOnce sync.Once

	ordersCreatedTotal metric.Int64Counter
	ordersPaidTotal    metric.Int64Counter
	ordersRevenueTotal metric.Int64Counter
)

// InitMetrics регистрирует бизнес-метрики. Вызывать после platform metrics.Init.
func InitMetrics() {
	metricsOnce.Do(func() {
		meter := otel.Meter("order-service")

		var err error

		ordersCreatedTotal, err = meter.Int64Counter("orders_created",
			metric.WithDescription("Количество созданных заказов"),
		)
		if err != nil {
			slog.Error("не удалось создать метрику orders_created", slog.String("error", err.Error()))
			return
		}

		ordersPaidTotal, err = meter.Int64Counter("orders_paid",
			metric.WithDescription("Количество оплаченных заказов"),
		)
		if err != nil {
			slog.Error("не удалось создать метрику orders_paid", slog.String("error", err.Error()))
			return
		}

		ordersRevenueTotal, err = meter.Int64Counter("orders_revenue",
			metric.WithDescription("Суммарная выручка (в копейках)"),
		)
		if err != nil {
			slog.Error("не удалось создать метрику orders_revenue", slog.String("error", err.Error()))
			return
		}
	})
}

func recordOrderCreated(ctx context.Context, orderUUID, userUUID uuid.UUID, totalPrice int64) {
	slog.InfoContext(ctx, "заказ создан",
		slog.String("order_uuid", orderUUID.String()),
		slog.String("user_uuid", userUUID.String()),
		slog.Int64("total_price", totalPrice),
	)
	ordersCreatedTotal.Add(ctx, 1)
}
