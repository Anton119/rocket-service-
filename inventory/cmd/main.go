package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	invapi "github.com/Anton119/rocket-service-/inventory/internal/api/inventory/v1"
	partrepo "github.com/Anton119/rocket-service-/inventory/internal/repository/part"
	partsvc "github.com/Anton119/rocket-service-/inventory/internal/service/part"
	"github.com/Anton119/rocket-service-/shared/pkg/config"
	"github.com/Anton119/rocket-service-/shared/pkg/grpc/interceptor"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

const (
	grpcAddress = "localhost:50051"

	shutdownTimeout = 10 * time.Second

	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	dsn, err := config.DBURI()
	if err != nil {
		slog.Error("получение DSN", "error", err)
		return err
	}

	pool, txManager, err := openPostgres(ctx, dsn)
	if err != nil {
		return err
	}
	defer pool.Close()

	lis, err := (&net.ListenConfig{}).Listen(context.Background(), "tcp", grpcAddress)
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		return err
	}

	grpcServer, err := newGRPCServer()
	if err != nil {
		return err
	}

	repo := partrepo.New(pool, txManager)
	catalog := partsvc.NewService(repo)
	api := invapi.NewAPI(catalog)
	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)
	reflection.Register(grpcServer)

	slog.Info("запуск InventoryService", "адрес", grpcAddress)

	return serveUntilSignal(grpcServer, lis)
}

func openPostgres(ctx context.Context, dsn string) (*pgxpool.Pool, *manager.Manager, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		slog.Error("создание пула соединений", "error", err)
		return nil, nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		slog.Error("проверка соединения с БД", "error", err)
		return nil, nil, err
	}
	slog.Info("подключение к PostgreSQL установлено")

	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		pool.Close()
		slog.Error("создание transaction manager", "error", err)
		return nil, nil, err
	}

	return pool, txManager, nil
}

func newGRPCServer() (*grpc.Server, error) {
	pvUnary, err := interceptor.UnaryProtovalidateInterceptor()
	if err != nil {
		slog.Error("ошибка инициализации protovalidate", "error", err)
		return nil, err
	}

	return grpc.NewServer(
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
		grpc.ChainUnaryInterceptor(
			interceptor.RecoveryInterceptor(),
			pvUnary,
			interceptor.LoggerInterceptor(),
		),
	), nil
}

func serveUntilSignal(grpcServer *grpc.Server, lis net.Listener) error {
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
		return gracefulStopGRPC(grpcServer, serveErrCh, sig)
	case serveErr := <-serveErrCh:
		if serveErr != nil {
			slog.Error("ошибка запуска сервера", "error", serveErr)
			return serveErr
		}
	}

	return nil
}

func gracefulStopGRPC(grpcServer *grpc.Server, serveErrCh <-chan error, sig os.Signal) error {
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

	return nil
}
