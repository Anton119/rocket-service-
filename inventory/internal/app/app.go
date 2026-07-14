package app

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

	"github.com/Anton119/rocket-service-/inventory/internal/config"
	invinterceptor "github.com/Anton119/rocket-service-/inventory/internal/interceptor"
	"github.com/Anton119/rocket-service-/platform/pkg/closer"
	"github.com/Anton119/rocket-service-/platform/pkg/grpc/health"
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

// App управляет жизненным циклом InventoryService.
type App struct {
	cfg         *config.Config
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

// New создаёт и инициализирует приложение.
func New(ctx context.Context, cfg *config.Config) *App {
	a := &App{cfg: cfg}
	a.initDeps(ctx)

	return a
}

// Run запускает gRPC-сервер и выполняет graceful shutdown по сигналу ОС.
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)
	go func() { errCh <- a.runGRPCServer() }()

	var runErr error
	select {
	case runErr = <-errCh:
	case <-ctx.Done():
		slog.Info("получен сигнал завершения, начинаем graceful shutdown")
	}
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout) //nolint:g118 // отдельный контекст для shutdown после отмены signal ctx.
	defer shutdownCancel()

	if err := closer.CloseAll(shutdownCtx); err != nil {
		slog.Error("ошибка при завершении работы", "error", err)
		if runErr == nil {
			runErr = err
		}
	}

	return runErr
}

func (a *App) initDeps(ctx context.Context) {
	inits := []func(context.Context){
		a.initDI,
		a.initLogger,
		a.initListener,
		a.initGRPCServer,
	}

	for _, f := range inits {
		f(ctx)
	}
}

func (a *App) initDI(_ context.Context) {
	a.diContainer = &diContainer{cfg: a.cfg}
}

func (a *App) initLogger(_ context.Context) {
	logger.Init(a.cfg.Logger.Level)
}

func (a *App) initListener(_ context.Context) {
	listener, err := net.Listen("tcp", a.cfg.GRPC.Address()) //nolint:noctx // net.Listen не требует контекст.
	if err != nil {
		slog.Error("не удалось создать TCP-листенер", "error", err)
		os.Exit(1)
	}

	a.listener = listener
}

func (a *App) initGRPCServer(ctx context.Context) {
	pvUnary, err := interceptor.UnaryProtovalidateInterceptor()
	if err != nil {
		slog.Error("ошибка инициализации protovalidate", "error", err)
		os.Exit(1)
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
			invinterceptor.UnaryErrorInterceptor(),
		),
	)

	closer.Add("gRPC server", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)
	health.RegisterService(a.grpcServer)
	inventoryv1.RegisterInventoryServiceServer(a.grpcServer, a.diContainer.InventoryV1API(ctx))
}

func (a *App) runGRPCServer() error {
	slog.Info("gRPC-сервер запущен", "address", a.cfg.GRPC.Address())

	return a.grpcServer.Serve(a.listener)
}
