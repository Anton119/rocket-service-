package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	invapi "github.com/Anton119/rocket-service-/inventory/internal/api/inventory/v1"
	partrepo "github.com/Anton119/rocket-service-/inventory/internal/repository/part"
	partsvc "github.com/Anton119/rocket-service-/inventory/internal/service/part"
	"github.com/Anton119/rocket-service-/shared/pkg/grpc/interceptor"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

const (
	// Адрес сервера.
	grpcAddress = "localhost:50051"

	// Таймауты для graceful shutdown.
	shutdownTimeout = 10 * time.Second

	// gRPC keepalive параметры.
	grpcMaxConnectionIdle     = 15 * time.Minute // Закрыть idle-соединения (нет активных RPC)
	grpcMaxConnectionAge      = 30 * time.Minute // Принудительная ротация для балансировки
	grpcMaxConnectionAgeGrace = 5 * time.Second  // Время на завершение активных RPC
	grpcKeepaliveTime         = 5 * time.Minute  // Интервал ping'ов для обнаружения мёртвых соединений
	grpcKeepaliveTimeout      = 1 * time.Second  // Таймаут ожидания pong
	grpcMinPingInterval       = 5 * time.Minute  // Минимальный интервал ping'ов от клиента (защита от DoS)
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	lis, err := (&net.ListenConfig{}).Listen(context.Background(), "tcp", grpcAddress)
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		return err
	}
	// прото валидация для gprc
	pvUnary, err := interceptor.UnaryProtovalidateInterceptor()
	if err != nil {
		slog.Error("ошибка инициализации protovalidate", "error", err)
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcMinPingInterval,
			PermitWithoutStream: true,
		}),

		// Интерцепторы: отлов паник, Protovalidate, логирование.
		grpc.ChainUnaryInterceptor(
			interceptor.RecoveryInterceptor(),
			pvUnary,
			interceptor.LoggerInterceptor(),
		),
	)
	repo := partrepo.NewRepository(partrepo.SeedParts())
	catalog := partsvc.NewService(repo)
	api := invapi.NewAPI(catalog)
	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)

	// Включаем reflection для postman/grpcurl
	reflection.Register(grpcServer)

	slog.Info("запуск InventoryService", "адрес", grpcAddress)

	serveErrCh := make(chan error, 1)
	go func() {
		slog.Info("🚀 gRPC сервер запущен", "адрес", grpcAddress)
		serveErrCh <- grpcServer.Serve(lis)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case sig := <-quit:
		slog.Info("🛑 завершение работы gRPC сервера...", "сигнал", sig.String())

		stopped := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(stopped)
		}()

		timer := time.NewTimer(shutdownTimeout)
		defer timer.Stop()
		select {
		case <-stopped:
			slog.Info("✅ сервер остановлен")
		case <-timer.C:
			slog.Warn("⏳ таймаут graceful shutdown, принудительная остановка")
			grpcServer.Stop()
		}

		if serveErr := <-serveErrCh; serveErr != nil {
			slog.Error("ошибка работы сервера", "error", serveErr)
			return serveErr
		}
	case serveErr := <-serveErrCh:
		if serveErr != nil {
			slog.Error("ошибка запуска сервера", "error", serveErr)
			return serveErr
		}
	}

	return nil
}
