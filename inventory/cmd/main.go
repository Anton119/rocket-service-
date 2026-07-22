package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/Anton119/rocket-service-/inventory/internal/app"
	"github.com/Anton119/rocket-service-/inventory/internal/config"
)

func main() {
	_ = godotenv.Load("inventory.env")    //nolint:gosec // .env файл опционален — ошибка загрузки допустима.
	_ = godotenv.Load("../inventory.env") //nolint:gosec // .env файл опционален — ошибка загрузки допустима.

	config.MustLoad(config.ResolveConfigPath())

	application := app.New(context.Background())

	slog.Info("запуск InventoryService", "адрес", config.AppConfig().GRPC.Address())

	if err := application.Run(); err != nil {
		slog.Error("InventoryService завершился с ошибкой", "error", err)
		os.Exit(1)
	}
}
