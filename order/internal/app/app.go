package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Anton119/rocket-service-/order/internal/config"
	ordersvc "github.com/Anton119/rocket-service-/order/internal/service/order"
	"github.com/Anton119/rocket-service-/platform/pkg/closer"
	"github.com/Anton119/rocket-service-/platform/pkg/logger"
	"github.com/Anton119/rocket-service-/platform/pkg/metrics"
	"github.com/Anton119/rocket-service-/platform/pkg/tracing"
)

// App — жизненный цикл OrderService.
type App struct {
	diContainer *diContainer
	httpServer  *http.Server
}

// New создаёт и инициализирует приложение.
func New(ctx context.Context) *App {
	a := &App{}
	a.initDeps(ctx)
	return a
}

// Run запускает HTTP-сервер и Kafka-потребитель, обрабатывает graceful shutdown.
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a.startGracefulShutdown(ctx, cancel)

	errCh := make(chan error, 2)

	go func() {
		slog.Info("🚀 HTTP-сервер запущен", "address", config.AppConfig().HTTP.Address())
		listenErr := a.httpServer.ListenAndServe()
		if listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			errCh <- listenErr
		}
	}()

	go func() {
		slog.Info("Kafka-потребитель ShipAssembled запущен",
			"topic", config.AppConfig().ShipAssembledConsumer.TopicName(),
			"group", config.AppConfig().ShipAssembledConsumer.ConsumerGroupID(),
		)
		if consErr := a.diContainer.AssemblyConsumerService().RunConsumer(ctx); consErr != nil && ctx.Err() == nil {
			errCh <- fmt.Errorf("потребитель упал: %w", consErr)
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
	a.initLogger()
	a.initMetrics()
	a.initTracing(ctx)
	a.diContainer = &diContainer{}
	a.httpServer = a.diContainer.HTTPServer(ctx)
}

func (a *App) initTracing(ctx context.Context) {
	shutdown, err := tracing.InitTracer(ctx, config.AppConfig().TracingPlatformConfig())
	if err != nil {
		slog.Error("не удалось инициализировать трейсер", "error", err)
		os.Exit(1)
	}

	closer.Add("tracer", shutdown)
}

func (a *App) initMetrics() {
	metrics.Init(
		config.AppConfig().Otel.GetServiceName(),
		metrics.WithCollectorEndpoint(config.AppConfig().Otel.CollectorEndpoint()),
	)
	ordersvc.InitMetrics()

	closer.Add("metrics", func(_ context.Context) error {
		return metrics.Close()
	})
}

func (a *App) initLogger() {
	logger.Init(config.AppConfig().LoggerPlatformConfig())
	closer.Add("logger", func(_ context.Context) error {
		return logger.Close()
	})
}

func (a *App) startGracefulShutdown(ctx context.Context, cancel context.CancelFunc) {
	go func() { //nolint:gosec // G118: ctx уже отменён, context.Background нужен для graceful shutdown.
		<-ctx.Done()
		cancel()

		slog.Info("получен сигнал завершения, начинаем graceful shutdown")

		shutdownCtx, shutdownCancel := context.WithTimeout(
			context.Background(),
			config.AppConfig().HTTP.ShutdownTimeout(),
		)
		defer shutdownCancel()

		if err := closer.CloseAll(shutdownCtx); err != nil {
			slog.Error("ошибка при завершении работы", "error", err)
		}
	}()
}
