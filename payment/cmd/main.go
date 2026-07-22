package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/Anton119/rocket-service-/payment/internal/app"
	"github.com/Anton119/rocket-service-/payment/internal/config"
)

func main() {
	_ = godotenv.Load("payment.env")    //nolint:gosec // .env файл опционален — ошибка загрузки допустима.
	_ = godotenv.Load("../payment.env") //nolint:gosec // .env файл опционален — ошибка загрузки допустима.

	config.MustLoad(config.ResolveConfigPath())

	application := app.New(context.Background())

	slog.Info("запуск PaymentService", "адрес", config.AppConfig().GRPC.Address())

	if err := application.Run(); err != nil {
		slog.Error("PaymentService завершился с ошибкой", "error", err)
		os.Exit(1)
	}
}
