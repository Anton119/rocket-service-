package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	orderapi "github.com/Anton119/rocket-service-/order/internal/api/order/v1"
	"github.com/Anton119/rocket-service-/order/internal/config"
	"github.com/Anton119/rocket-service-/platform/pkg/closer"
	"github.com/Anton119/rocket-service-/platform/pkg/logger"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	shutdownTimeout   = 10 * time.Second
)

// App управляет жизненным циклом OrderService.
type App struct {
	cfg         *config.Config
	diContainer *diContainer
	httpServer  *http.Server
}

// New создаёт и инициализирует приложение.
func New(ctx context.Context, cfg *config.Config) *App {
	a := &App{cfg: cfg}
	a.initDeps(ctx)

	return a
}

// Run запускает HTTP-сервер и выполняет graceful shutdown по сигналу ОС.
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)
	go func() { errCh <- a.runHTTPServer() }()

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
		a.initHTTPServer,
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

func (a *App) initHTTPServer(ctx context.Context) {
	orderServer, err := orderapi.NewServer(a.diContainer.OrderAPI(ctx))
	if err != nil {
		slog.Error("ошибка создания HTTP-сервера OpenAPI", "error", err)
		os.Exit(1)
	}

	a.httpServer = &http.Server{
		Addr:              a.cfg.HTTP.Address(),
		Handler:           orderServer,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	closer.Add("HTTP server", func(shutdownCtx context.Context) error {
		return a.httpServer.Shutdown(shutdownCtx)
	})
}

func (a *App) runHTTPServer() error {
	slog.Info("HTTP-сервер запущен", "address", a.cfg.HTTP.Address())

	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
