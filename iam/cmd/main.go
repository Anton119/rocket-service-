package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/Anton119/rocket-service-/iam/internal/app"
	"github.com/Anton119/rocket-service-/iam/internal/config"
)

func main() {
	_ = godotenv.Load("iam.env")    //nolint:gosec // .env файл опционален — ошибка загрузки допустима.
	_ = godotenv.Load("../iam.env") //nolint:gosec // .env файл опционален — ошибка загрузки допустима.

	config.MustLoad(config.ResolveConfigPath())

	application := app.New(context.Background())

	slog.Info("запуск IAMService", "адрес", config.AppConfig().GRPC.Address())

	if err := application.Run(); err != nil {
		slog.Error("IAMService завершился с ошибкой", "error", err)
		os.Exit(1)
	}
}
