package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/Anton119/rocket-service-/assembly/internal/app"
	"github.com/Anton119/rocket-service-/assembly/internal/config"
)

func main() {
	_ = godotenv.Load("assembly.env")    //nolint:gosec // .env файл опционален — ошибка загрузки допустима.
	_ = godotenv.Load("../assembly.env") //nolint:gosec // .env файл опционален — ошибка загрузки допустима.

	config.MustLoad(config.ResolveConfigPath())

	application := app.New()

	slog.Info("запуск AssemblyService",
		"order_paid_topic", config.AppConfig().OrderPaidConsumer.TopicName(),
		"ship_assembled_topic", config.AppConfig().ShipAssembledProducer.TopicName(),
	)

	if err := application.Run(); err != nil {
		slog.Error("AssemblyService завершился с ошибкой", "error", err)
		os.Exit(1)
	}
}
