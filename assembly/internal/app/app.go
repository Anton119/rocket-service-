package app

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/Anton119/rocket-service-/assembly/internal/config"
	"github.com/Anton119/rocket-service-/platform/pkg/closer"
	"github.com/Anton119/rocket-service-/platform/pkg/logger"
)

const shutdownTimeout = 10 * time.Second

// App — жизненный цикл AssemblyService.
type App struct {
	diContainer *diContainer
}

// New создаёт и инициализирует приложение.
func New(_ context.Context) *App {
	a := &App{}
	a.initDeps()
	return a
}

// Run запускает consumer и обрабатывает graceful shutdown.
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a.startGracefulShutdown(ctx, cancel)

	errCh := make(chan error, 1)
	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("потребитель упал: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

func (a *App) initDeps() {
	a.initLogger()
	a.diContainer = &diContainer{}
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

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()

		if err := closer.CloseAll(shutdownCtx); err != nil {
			slog.Error("ошибка при завершении работы", "error", err)
		}
	}()
}

func (a *App) runConsumer(ctx context.Context) error {
	slog.Info("Kafka-потребитель OrderPaid запущен",
		"topic", config.AppConfig().OrderPaidConsumer.TopicName(),
		"group", config.AppConfig().OrderPaidConsumer.ConsumerGroupID(),
		"brokers", config.AppConfig().Kafka.Brokers,
	)
	return a.diContainer.OrderPaidConsumerService().RunConsumer(ctx)
}
