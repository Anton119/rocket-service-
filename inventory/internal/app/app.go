package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/Anton119/rocket-service-/inventory/internal/config"
	invinterceptor "github.com/Anton119/rocket-service-/inventory/internal/interceptor"
	"github.com/Anton119/rocket-service-/platform/pkg/closer"
	"github.com/Anton119/rocket-service-/platform/pkg/logger"
	"github.com/Anton119/rocket-service-/shared/pkg/grpc/interceptor"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

const (
	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute
	shutdownTimeout           = 10 * time.Second
)

// App — жизненный цикл InventoryService.
type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

// New создаёт и инициализирует приложение.
func New(ctx context.Context) *App {
	a := &App{}
	a.initDeps(ctx)
	return a
}

// Run запускает gRPC-сервер и обрабатывает graceful shutdown.
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a.startGracefulShutdown(ctx, cancel)

	errCh := make(chan error, 1)
	go func() {
		slog.Info("🚀 gRPC сервер запущен", "адрес", config.AppConfig().GRPC.Address())
		if err := a.grpcServer.Serve(a.listener); err != nil {
			errCh <- fmt.Errorf("gRPC-сервер упал: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

func (a *App) initDeps(ctx context.Context) {
	a.diContainer = &diContainer{}
	logger.Init(config.AppConfig().Logger.Level)
	a.initListener(ctx)
	a.initGRPCServer(ctx)
}

func (a *App) initListener(_ context.Context) {
	listener, err := (&net.ListenConfig{}).Listen(context.Background(), "tcp", config.AppConfig().GRPC.Address())
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		panic(err)
	}
	a.listener = listener
}

func (a *App) initGRPCServer(ctx context.Context) {
	pvUnary, err := interceptor.UnaryProtovalidateInterceptor()
	if err != nil {
		slog.Error("ошибка инициализации protovalidate", "error", err)
		panic(err)
	}

	a.grpcServer = grpc.NewServer(
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
			invinterceptor.AuthInterceptor(a.diContainer.AuthClient()),
			invinterceptor.ErrorInterceptor(),
		),
	)

	closer.Add("gRPC server", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	inventoryv1.RegisterInventoryServiceServer(a.grpcServer, a.diContainer.InventoryAPI(ctx))
	reflection.Register(a.grpcServer)
}

func (a *App) startGracefulShutdown(ctx context.Context, cancel context.CancelFunc) {
	go func() { //nolint:gosec // G118: ctx уже отменён, context.Background нужен для graceful shutdown.
		<-ctx.Done()
		cancel()

		slog.Info("🛑 завершение работы gRPC сервера...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()

		if err := closer.CloseAll(shutdownCtx); err != nil {
			slog.Error("ошибка при завершении работы", "error", err)
		}
	}()
}
