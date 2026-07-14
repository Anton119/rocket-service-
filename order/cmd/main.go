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
	_ = godotenv.Load("order.env") //nolint:gosec // .env опционален — ошибка загрузки допустима.

	configPath := config.ResolveConfigPath()

	cfg, err := config.Load(configPath)
	if err != nil {
		slog.Error("не удалось загрузить конфигурацию", "error", err, "config_path", configPath)
		os.Exit(1)
	}

	if err = app.New(context.Background(), cfg).Run(); err != nil {
		slog.Error("ошибка запуска приложения", "error", err)
		os.Exit(1)
	}
}
