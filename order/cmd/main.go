package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderHandler "github.com/Anton119/rocket-service-/order/pkg/handler"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/payment/v1"
)

const (
	inventoryServiceAddress = "localhost:50051"
	paymentServiceAddress   = "localhost:50052"
)

const (
	httpPort     = "8080"
	urlParamCity = "city"

	// Таймауты для HTTP-сервера.
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	shutdownTimeout   = 10 * time.Second
	middlewareTimeout = 10 * time.Second
)

func main() {
	if err := run(); err != nil {
		os.Exit(1) // тут уже можно: все defer внутри run уже отработали
	}
}

func run() error {
	inventoryConn, err := grpc.NewClient(inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("не удалось подключиться к InventoryService", "error", err)
		return err
	}
	defer inventoryConn.Close()

	paymentConn, err := grpc.NewClient(paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("не удалось подключиться к PaymentService", "error", err)
		return err
	}
	defer paymentConn.Close()

	// Создаём хранилище и обработчик
	store := orderHandler.NewOrderStore()
	h := orderHandler.NewOrderHandler(
		inventoryv1.NewInventoryServiceClient(inventoryConn),
		paymentv1.NewPaymentServiceClient(paymentConn),
		store,
	)

	// Создать OpenAPI сервер
	orderServer, err := orderHandler.SetupServer(h)
	if err != nil {
		slog.Error("ошибка создания сервера OpenAPI", "error", err)
		return err
	}

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           orderServer,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атаки
		ReadTimeout:       readTimeout,       // Лимит на чтение всего запроса
		WriteTimeout:      writeTimeout,      // Лимит на запись ответа
		IdleTimeout:       idleTimeout,       // Таймаут keep-alive соединений

	}

	go func() {
		slog.Info("🚀 HTTP-сервер запущен на порту", "port", httpPort)
		listenErr := server.ListenAndServe()
		if listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			slog.Error("❌ ошибка запуска сервера", "error", listenErr)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("🛑 завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("❌ ошибка при остановке сервера", "error", err)
	}
	slog.Info("✅ сервер остановлен")

	return nil
}
