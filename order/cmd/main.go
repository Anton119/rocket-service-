package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/Anton119/rocket-service-/order/internal/app"
	"github.com/Anton119/rocket-service-/order/internal/config"
)

func main() {
	_ = godotenv.Load("order.env")    //nolint:gosec // .env файл опционален — ошибка загрузки допустима.
	_ = godotenv.Load("../order.env") //nolint:gosec // .env файл опционален — ошибка загрузки допустима.

	config.MustLoad(config.ResolveConfigPath())

	application := app.New(context.Background())

	slog.Info("запуск OrderService",
		"http", config.AppConfig().HTTP.Address(),
		"order_paid_topic", config.AppConfig().OrderPaidProducer.TopicName(),
		"ship_assembled_topic", config.AppConfig().ShipAssembledConsumer.TopicName(),
	)

	if err := application.Run(); err != nil {
		slog.Error("OrderService завершился с ошибкой", "error", err)
		os.Exit(1)
	}
}
